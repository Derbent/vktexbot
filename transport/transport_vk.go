package transport

import (
	"context"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	longpoll "github.com/SevereCloud/vksdk/v2/longpoll-bot"
)

type vkTransport struct {
	vk *api.VK
	lp *longpoll.LongPoll
}

func NewVK(token string, groupID int) Transport {
	vk := api.NewVK(token)
	lp, err := longpoll.NewLongPoll(vk, groupID)

	if err != nil {
		panic(err)
	}

	return &vkTransport{
		vk: vk,
		lp: lp,
	}
}

func (vt *vkTransport) MessageNew(f func(*Client)) {
	vt.lp.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		f(&Client{
			id: obj.Message.PeerID,
			message: Message{
				text: &obj.Message.Text,
			},
		})
	})
}

func (vt *vkTransport) MessageSend(client *Client) {
	b := params.NewMessagesSendBuilder()
	b.Message(*client.message.text)
	b.RandomID(0)
	b.PeerID(client.id)
	vt.vk.MessagesSend(b.Params)
}

func (vt *vkTransport) Run() {
	vt.lp.Run()
}
