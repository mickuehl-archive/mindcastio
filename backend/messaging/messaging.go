package messaging

import (
	"github.com/nats-io/nats"

	"github.com/mindcastio/mindcastio/backend/environment"
	"github.com/mindcastio/mindcastio/backend/logger"
)

type MessageBroker struct {
	connection *nats.Conn
}

var broker *MessageBroker

func Initialize(env *environment.Environment) {

	logger.Log("messaging.initialize", env.MessagingServiceUrls()[0])

	br := MessageBroker{}

	opts := nats.DefaultOptions
	opts.Servers = env.MessagingServiceUrls()

	nc, err := opts.Connect()
	if err != nil {
		panic(err)
	}
	br.connection = nc

	// make it *global*
	broker = &br
}

func Shutdown() {
	logger.Log("messaging.shutdown")

	broker.connection.Close()
}

func Send(subj string, msg string) error {
	return broker.connection.Publish(subj, []byte(msg))
}

func Subscribe(subj string, callback func(msg *nats.Msg)) {
	broker.connection.Subscribe(subj, callback)
}

func QueueSubscribe(subj string, queue string, callback func(msg *nats.Msg)) {
	broker.connection.QueueSubscribe(subj, queue, callback)
}
