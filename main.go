package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

// Car ...
type Car struct {
	Mark    string `json:"mark"`
	Model   string `json:"model"`
	Age     int    `json:"age"`
	Carcase string `json:"carcase"`
}

var messageChan = make(chan string)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")

		defer func() {
			close(messageChan)
			messageChan = nil
		}()

		flusher, ok := w.(http.Flusher)
		if !ok {
			log.Println("not found http.Flusher")
		}

		for {
			select {
			case message := <-messageChan:
				fmt.Fprintf(w, "data: %s\n\n", message)
				flusher.Flush()
			case <-r.Context().Done():
				return
			}
		}

	})

	mux.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if messageChan != nil {
			cars := make([]Car, 0)

			for {
				cars = append(cars, Car{
					Mark:    "Mercedes-Benz",
					Model:   "e63",
					Age:     2015,
					Carcase: "SUV",
				})

				bCars, err := json.Marshal(cars)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				log.Printf("print message to client")
				message := string(bCars)
				// send the message through the available channel
				messageChan <- message
				time.Sleep(10 * time.Second)
			}
		}

	})
	http.ListenAndServe(":4500", mux)
}
