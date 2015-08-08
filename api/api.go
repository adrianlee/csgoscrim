package main

import (
	"fmt"
	"github.com/pusher/pusher-http-go"
	"net/http"
)

func main() {
	client := pusher.Client{
		AppId:  "133024",
		Key:    "ba691b85219d373cdeeb",
		Secret: "1a17090c04e6ed9d06f5",
	}

	data := map[string]string{"message": "hello world"}

	client.Trigger("test_channel", "my_event", data)

	// http.HandleFunc("/", handler)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, world!")
}
