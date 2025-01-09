package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

type GyroscopeData struct {
	AngleX float64
	AngleY float64
	AngleZ float64
}

type Config struct {
	X int
	Y int
	Z int
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

	// Get a new midi driver port
	port, err := rtmididrv.New()
	if err != nil {
		panic(err)
	}

	// Open the port as a virtual output
	out, err := port.OpenVirtualOut("GyroMidi")
	if err != nil {
		panic(err)
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

			log.Printf("Received gyroscope data: (X: %v, Y: %v, Z: %v)", gyro.AngleX, gyro.AngleY, gyro.AngleZ)

			var cc_x uint8 = uint8((127.0 / 360.0) * gyro.AngleX)
			var cc_y uint8 = uint8((127.0 / 360.0) * gyro.AngleY)
			var cc_z uint8 = uint8((127.0 / 360.0) * gyro.AngleZ)

			p := uint8(out.Number())
			out.Send(midi.ControlChange(p, uint8(cfg.X), cc_x))
			out.Send(midi.ControlChange(p, uint8(cfg.Y), cc_y))
			out.Send(midi.ControlChange(p, uint8(cfg.Z), cc_z))
		}

	}

	// Routing
	http.HandleFunc("/", controller_handler)
	http.HandleFunc("/action", action_handler)

	log.Printf("Hosting on https://%v:8080/\n", GetOutboundIP())
	log.Fatal(http.ListenAndServeTLS(":8080", "ssl/gyromidi.crt", "ssl/gyromidi.key", nil)) // Start server with ssl
}
