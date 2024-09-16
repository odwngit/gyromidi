package main

import (
	"log"
	"encoding/json"
	"net/http"
	"github.com/go-vgo/robotgo"
	"github.com/BurntSushi/toml"
)

type GyroscopeData struct {
	AngleX float64
	AngleY float64
	AngleZ float64
}

type Config struct {
	Mode string
	Axis string
	
	Sensitivity float64
	AngleOffset float64

	ThresholdLeft float64
	ThresholdRight float64

	ThresholdLeftKey string
	ThresholdRightKey string
}

func main() {
	log.Println("Starting starwheel...")

	var cfg Config // Load config.toml
	_, err := toml.DecodeFile("config.toml", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Config loaded from config.toml: Mode = '%v'\n", cfg.Mode)

	sx, sy := robotgo.GetScreenSize()
	var left_pressed bool = false
	var right_pressed bool = false
	var gyro GyroscopeData

	controller_handler := func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "site/controller.html") // Serve up file
	}

	interface_handler := func(writer http.ResponseWriter, request *http.Request) {
		http.ServeFile(writer, request, "site/interface.html") // Serve up file
	}

	action_handler := func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "GET" {
			writer.Header().Set("Content-Type", "application/json") // Write json header
			writer.WriteHeader(http.StatusCreated) // Send http status
			json.NewEncoder(writer).Encode(gyro) // Send back config data

		} else if request.Method == "POST" {
			decoder := json.NewDecoder(request.Body) // Decode data sent from POST request
			err := decoder.Decode(&gyro) // And decode it into GyroscopeData
			if err != nil {
				panic(err)
			}
			
			// Get selected axis
			var selected_axis float64

			switch cfg.Axis {
				case "X":
					selected_axis = gyro.AngleX
				case "Y":
					selected_axis = gyro.AngleY
				case "Z":
					selected_axis = gyro.AngleZ
				default:
					log.Fatal("Error: unsupported config value: Axis")
			}

			selected_axis += cfg.AngleOffset

			// Handle gyro data
			switch cfg.Mode {
				case "mouse":
					robotgo.Move((sx/2)+(int(selected_axis*cfg.Sensitivity)), sy/2)
				case "threshold":
					if selected_axis < cfg.ThresholdLeft && !left_pressed {
						left_pressed = true
						robotgo.KeyToggle(cfg.ThresholdLeftKey)
					} else if selected_axis > cfg.ThresholdRight && !right_pressed {
						right_pressed = true
						robotgo.KeyToggle(cfg.ThresholdRightKey)
					} else if (left_pressed || right_pressed) && selected_axis > cfg.ThresholdLeft && selected_axis < cfg.ThresholdRight {
						left_pressed = false
						right_pressed = false
						robotgo.KeyToggle(cfg.ThresholdLeftKey, "up")
						robotgo.KeyToggle(cfg.ThresholdRightKey, "up")
					}
				default:
					log.Fatal("Error: unsupported config value: Mode")
			}

			selected_axis -= cfg.AngleOffset

		}
	}

	config_handler := func(writer http.ResponseWriter, request *http.Request) {
		if request.Method == "GET" {
			writer.Header().Set("Content-Type", "application/json") // Write json header
			writer.WriteHeader(http.StatusCreated) // Send http status
			json.NewEncoder(writer).Encode(cfg) // Send back config data
		} else if request.Method == "POST" {
			decoder := json.NewDecoder(request.Body) // Decode data sent from POST request
			err := decoder.Decode(&cfg) // And decode it into Config
			if err != nil {
				panic(err)
			}
		}
	}

	http.HandleFunc("/starwheel", controller_handler) // Route /starwheel into controller_handler
	http.HandleFunc("/config", config_handler) // Route /config into config_handler
	http.HandleFunc("/action", action_handler) // Route /action into action_handler
	http.HandleFunc("/", interface_handler) // Route / into interface_handler
	log.Printf("Hosting on https://%v:8080/\n", GetOutboundIP())

	log.Fatal(http.ListenAndServeTLS(":8080", "ssl/starwheel.crt", "ssl/starwheel.key", nil)) // Start server with ssl
}
