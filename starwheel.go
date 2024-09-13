package main

import (
	"fmt"
	"log"
	"encoding/json"
	"net/http"
	//"github.com/go-vgo/robotgo"
)

type GyroscopeData struct {
	AngleX float64
	AngleY float64
	AngleZ float64
}

func main() {
	fmt.Println("Starting http server...")
	var gyro GyroscopeData

	siteHandler := func(writer http.ResponseWriter, request *http.Request) { // Site handler
		if request.Method == "POST" { // If receiving a POST request
			decoder := json.NewDecoder(request.Body) // Decode data sent from POST request
			err := decoder.Decode(&gyro) // And decode it into GyroscopeData
			if err != nil {
				panic(err)
			}
			fmt.Println(gyro.AngleY)
		} else {
			http.ServeFile(writer, request, "controller/index.html") // Serve up file
		}
	}

	http.HandleFunc("/starwheel", siteHandler)
	
	fmt.Printf("Hosting on %v:8080/starwheel\n", GetOutboundIP())
	log.Fatal(http.ListenAndServe(":8080", nil))
}
