package transport

import (
	"context"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/events"
	longpoll "github.com/SevereCloud/vksdk/v2/longpoll-bot"
	"github.com/pkg/errors"
)

type vkTransport struct {
	vk *api.VK
	lp *longpoll.LongPoll
}

func NewVK(token string, groupID int) (Transport, error) {
	vk := api.NewVK(token)
	lp, err := longpoll.NewLongPoll(vk, groupID)

	if err != nil {
		return nil, err
	}

	return &vkTransport{
		vk: vk,
		lp: lp,
	}, nil
}

func (vt *vkTransport) MessageNew(f func(*Client)) {
	vt.lp.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		f(&Client{
			id: obj.Message.PeerID,
			message: Message{
				text: obj.Message.Text,
			},
		})
	})
}

func (vt *vkTransport) MessageSend(client *Client) error {
	b := params.NewMessagesSendBuilder()
	b.Message(client.message.text)
	b.RandomID(0)
	b.PeerID(client.id)
	_, err := vt.vk.MessagesSend(b.Params)
	return errors.Wrap(err, "MessageSend error")
}

func (vt *vkTransport) Run() {
	vt.lp.Run()
}
