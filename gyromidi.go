package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/BurntSushi/toml"
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

			// Do midi stuff
		}

	}

	// Routing
	http.HandleFunc("/", controller_handler)
	http.HandleFunc("/action", action_handler)

	log.Printf("Hosting on https://%v:8080/\n", GetOutboundIP())
	log.Fatal(http.ListenAndServeTLS(":8080", "ssl/gyromidi.crt", "ssl/gyromidi.key", nil)) // Start server with ssl
}
