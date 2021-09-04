package transport

type Message struct {
	text *string
}

type Client struct {
	id      int
	message Message
}

type Transport interface {
	MessageNew(func(*Client))
	MessageSend(*Client)
	Run()
}

func Run(t Transport) {
	t.MessageNew(func(client *Client) {
		t.MessageSend(client)
	})
	t.Run()
}
