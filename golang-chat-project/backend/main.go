package main

import (
	"fmt"
	"net/http"

	"github.com/dev-dhanushkumar/golang-chat/pkg/websocket"
)

func serveWS(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	fmt.Println("Websocket endpoint reached!")

	conn, err := websocket.Upgrade(w, r)

	if err != nil {
		fmt.Println(w, "%+v\n", err)
	}

	client := &websocket.Client{
		Conn: conn,
		Pool: pool,
	}

	pool.Register <- client
	client.Read()
}

func setupRoute() {
	pool := websocket.NewPool()
	go pool.Start()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(pool, w, r)
	})
}

func main() {
	fmt.Println("Dhanush's full stack chat Project")
	setupRoute()
	http.ListenAndServe(":9000", nil)
}
