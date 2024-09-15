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
	
	Sensitivity float64

	ThresholdLeft float64
	ThresholdRight float64

	ThresholdLeftKey string
	ThresholdRightKey string
}

func main() {
	log.Println("Starting http server...")

	var gyro GyroscopeData

	var cfg Config
	_, err := toml.DecodeFile("config.toml", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Config loaded from config.toml")

	sx, sy := robotgo.GetScreenSize()

	var left_pressed bool = false
	var right_pressed bool = false

	site_handler := func(writer http.ResponseWriter, request *http.Request) { // Site handler
		if request.Method == "POST" { // If receiving a POST request
			decoder := json.NewDecoder(request.Body) // Decode data sent from POST request
			err := decoder.Decode(&gyro) // And decode it into GyroscopeData
			if err != nil {
				panic(err)
			}

			// Handle gyro data
			switch cfg.Mode {
				case "mouse":
					robotgo.Move((sx/2)+(int(gyro.AngleY*cfg.Sensitivity)), sy/2)
				case "threshold":
					if gyro.AngleY < cfg.ThresholdLeft && !left_pressed {
						left_pressed = true
						robotgo.KeyToggle(cfg.ThresholdLeftKey)
					} else if gyro.AngleY > cfg.ThresholdRight && !right_pressed {
						right_pressed = true
						robotgo.KeyToggle(cfg.ThresholdRightKey)
					} else if (left_pressed || right_pressed) && gyro.AngleY > cfg.ThresholdLeft && gyro.AngleY < cfg.ThresholdRight {
						left_pressed = false
						right_pressed = false
						robotgo.KeyToggle(cfg.ThresholdLeftKey, "up")
						robotgo.KeyToggle(cfg.ThresholdRightKey, "up")
					}
				default:
					log.Fatal("Error: unsupported config Mode")
			}

		} else {
			http.ServeFile(writer, request, "controller/index.html") // Serve up file
		}
	}

	http.HandleFunc("/starwheel", site_handler) // Route /starwheel into site_handler
	log.Printf("Hosting on %v:8080/starwheel\n", GetOutboundIP())

	log.Fatal(http.ListenAndServeTLS(":8080", "ssl/localhost.crt", "ssl/localhost.key", nil)) // Start server with ssl
}
