package main

import (
	"log"
	"encoding/json"
	"net/http"
	"github.com/go-vgo/robotgo"
)

type GyroscopeData struct {
	AngleX float64
	AngleY float64
	AngleZ float64
}

func main() {
	log.Println("Starting http server...")
	var gyro GyroscopeData

	sx, sy := robotgo.GetScreenSize()

	siteHandler := func(writer http.ResponseWriter, request *http.Request) { // Site handler
		if request.Method == "POST" { // If receiving a POST request
			decoder := json.NewDecoder(request.Body) // Decode data sent from POST request
			err := decoder.Decode(&gyro) // And decode it into GyroscopeData
			if err != nil {
				panic(err)
			}

			var mousex int = (sx/2)+(int(gyro.AngleY)*6)
			robotgo.Move(mousex, sy/2)

		} else {
			http.ServeFile(writer, request, "controller/index.html") // Serve up file
		}
	}

	http.HandleFunc("/starwheel", siteHandler)
	
	log.Printf("Hosting on %v:8080/starwheel\n", GetOutboundIP())
	log.Fatal(http.ListenAndServeTLS(":8080", "ssl/localhost.crt", "ssl/localhost.key", nil))
}
