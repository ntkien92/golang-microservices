package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Time struct {
	CurrentTime string `json:"current_time"`
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		currentTime := []Time{
			{CurrentTime: time.Now().Format(http.TimeFormat)},
		}
		json.NewEncoder(w).Encode(currentTime)
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
