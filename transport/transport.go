package transport

import (
	"fmt"
	"image"
	"log"

	"github.com/ZashX/vktexbot/render"
	"github.com/pkg/errors"
)

type Message struct {
	text  string
	image *image.RGBA
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
		log.Printf("Message is recieved from \"%v\", text \"%v\"", client.id, client.message.text)

		if len(client.message.text) > 0 {
			r := render.New()
			image, err := r.Rend(client.message.text)

			if err != nil {
				log.Printf("%v, text \"%v\"", errors.Wrap(err, "Error compilation text"), client.message.text)
				client.message = Message{
					text:  fmt.Sprint(err),
					image: nil,
				}
			} else {
				client.message = Message{
					text:  "",
					image: image,
				}
			}

			if err := t.MessageSend(client); err != nil {
				log.Print(err)
			} else {
				log.Printf("Message succesfuly sended to %v", client.id)
			}
		}
	})
	t.Run()
}
