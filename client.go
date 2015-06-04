package main

import (
	"github.com/gorilla/websocket"
	"time"
	"log"
)

type client struct {
	socket *websocket.Conn
	send chan *message
	room *room
	userData map[string]interface {}
}

func (c *client) read() {
	for {
		var msg *message
		log.Println(c.userData)
		if err := c.socket.ReadJSON(&msg); err == nil {
			msg.When = time.Now().Format("15:04")
			msg.Name = c.userData["name"].(string)
			if avatarURL, ok := c.userData["avatarURL"]; ok {
				msg.AvatarURL = avatarURL.(string)
			}
			c.room.forward <- msg
		} else {
			break;
		}
	}
	c.socket.Close()
}

func (c *client) write() {
	for msg := range c.send {
		if err := c.socket.WriteJSON(msg); err != nil {
			break;
		}
	}
	c.socket.Close()
}
