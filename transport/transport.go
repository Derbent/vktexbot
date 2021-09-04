package transport

import (
	"log"
)

type Message struct {
	text string
}

type Client struct {
	id      int
	message Message
}

type Transport interface {
	MessageNew(func(*Client))
	MessageSend(*Client) error
	Run()
}

func Run(t Transport) {
	t.MessageNew(func(client *Client) {
		log.Printf("message is recieved from %v, text \"%v\"", client.id, client.message.text)
		if len(client.message.text) > 0 {
			if err := t.MessageSend(client); err != nil {
				log.Print(err)
			}
		}
	})
	t.Run()
}
