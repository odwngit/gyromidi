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

	site_handler := func(writer http.ResponseWriter, request *http.Request) { // Site handler
		if request.Method == "POST" { // If receiving a POST request
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

		} else {
			http.ServeFile(writer, request, "site/controller.html") // Serve up file
		}
	}

	http.HandleFunc("/starwheel", site_handler) // Route /starwheel into site_handler
	log.Printf("Hosting on https://%v:8080/starwheel\n", GetOutboundIP())

	log.Fatal(http.ListenAndServeTLS(":8080", "ssl/starwheel.crt", "ssl/starwheel.key", nil)) // Start server with ssl
}
