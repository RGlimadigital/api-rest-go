package models

import (
	"io/ioutil"
	"log"
	"net/http"
	"github.com/gorilla/websocket"

)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var wsConnections = []*websocket.Conn{}

func ConnectWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error creating websocket connection", err)
		w.WriteHeader(http.StatusInternalServerError)
	}
	wsConnections = append(wsConnections, conn)
	for {
		messageType, reader, err := conn.NextReader()
		if err != nil {
			log.Println("Error processing message", err)
			break
		}
		messageBytes, err := ioutil.ReadAll(reader)
		println("Nuevo Mensaje: " + string(messageBytes))
		if err == nil {
			for _, otherConn := range wsConnections {
				if otherConn != conn {
					otherConn.WriteMessage(messageType, messageBytes)
				}
			}
		}
	}
	conn.Close()
}
