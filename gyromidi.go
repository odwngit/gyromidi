package main

import (
	"encoding/json"
	"log"
	"net/http"
	"runtime"

	"github.com/BurntSushi/toml"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type GyroscopeData struct {
	AngleX       float64
	AngleY       float64
	AngleZ       float64
	Acceleration float64
}

type Config struct {
	X              int
	Y              int
	Z              int
	A              int
	VerboseLogging bool
}

func main() {
	log.Println("Starting gyromidi...")

	var cfg Config // Load config.toml
	_, err := toml.DecodeFile("config.toml", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Config loaded from config.toml: (CC %v, %v, %v)", cfg.X, cfg.Y, cfg.Z)

	var gyro GyroscopeData

	var out drivers.Out

	switch runtime.GOOS {
	case "darwin", "linux":
		// Get a new midi driver
		port, err := rtmididrv.New()
		if err != nil {
			log.Fatal(err)
		}

		// Open the driver as a virtual output port
		out, err = port.OpenVirtualOut("GyroMidi")
		if err != nil {
			log.Fatal(err)
		}
	case "windows":
		out, err = drivers.OutByName("GyroMidi")
		if err != nil {
			log.Println("An error occurred. Have you set up loopMIDI correctly?")
			log.Fatal(err)
		}
	default:
		log.Fatal("Your platform is not supported.")
	}

	log.Printf("Opened virtual midi output %v...", out)

	controller_handler := func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "site/controller.html")
	}

	action_handler := func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "GET" {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(http.StatusAccepted)
			json.NewEncoder(writer).Encode(gyro)
		} else if request.Method == "POST" {
			decoder := json.NewDecoder(request.Body) // Decode data sent from POST request
			err := decoder.Decode(&gyro)             // And decode it into GyroscopeData
			if err != nil {
				panic(err)
			}

			if cfg.VerboseLogging {
				log.Printf("Received motion data: (X: %v, Y: %v, Z: %v, Acc: %v)", gyro.AngleX, gyro.AngleY, gyro.AngleZ, gyro.Acceleration)
			}

			// This should definitely have some sanitisation at some point
			var cc_x uint8 = uint8((127.0 / 360.0) * gyro.AngleX)
			var cc_y uint8 = uint8((127.0 / 360.0) * gyro.AngleY)
			var cc_z uint8 = uint8((127.0 / 360.0) * gyro.AngleZ)
			var cc_a uint8 = uint8(gyro.Acceleration)

			p := uint8(out.Number())

			if cfg.X != 0 {
				out.Send(midi.ControlChange(p, uint8(cfg.X), cc_x))
			}
			if cfg.Y != 0 {
				out.Send(midi.ControlChange(p, uint8(cfg.Y), cc_y))
			}
			if cfg.Z != 0 {
				out.Send(midi.ControlChange(p, uint8(cfg.Z), cc_z))
			}
			if cfg.A != 0 {
				out.Send(midi.ControlChange(p, uint8(cfg.A), cc_a))
			}
		}

	}

	// Routing
	http.HandleFunc("/", controller_handler)
	http.HandleFunc("/action", action_handler)

	log.Printf("Hosting on https://%v:8080/\n", GetOutboundIP())
	log.Fatal(http.ListenAndServeTLS(":8080", "ssl/gyromidi.crt", "ssl/gyromidi.key", nil)) // Start server with ssl
}
