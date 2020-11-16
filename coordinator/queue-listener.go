package coordinator

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"my-go-apps/InventoryService/dto"
	queueutils "my-go-apps/InventoryService/queue-utils"

	"github.com/streadway/amqp"
)

const url = "amqp://guest:guest@localhost:5672"

type QueueListener struct {
	conn    *amqp.Connection
	ch      *amqp.Channel
	sources map[string]<-chan amqp.Delivery
	ea      *EventAggregator
}

func NewQueueListener(ea *EventAggregator) *QueueListener {
	ql := QueueListener{
		sources: make(map[string]<-chan amqp.Delivery), //a channel is a data structure
		ea:      ea,
	}

	ql.conn, ql.ch = queueutils.GetChannel(url)

	return &ql
}

func (ql *QueueListener) ListenForNewSource() {
	q := queueutils.GetQueue("", ql.ch, true)
	//by default queues are set up with the default exchange, so need to rebind to associate it with a fanout exchange
	ql.ch.QueueBind(
		q.Name,
		"",
		"",
		false,
		nil,
	)

	//once the queue is bound, we can consume the messages queued onto it
	msgs, _ := ql.ch.Consume( //wire up a receiever for the messages that are sent to it
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	ql.DiscoverInventories()

	fmt.Println("listening for new sources")

	for msg := range msgs {
		//check if the new message has already been registered
		if ql.sources[string(msg.Body)] == nil {
			ql.ea.PublishEvent("DataSourceDiscovered", string(msg.Body))
			sourceChan, _ := ql.ch.Consume(
				"",
				"",
				true,
				false,
				false,
				false,
				nil,
			)

			ql.sources[string(msg.Body)] = sourceChan
			go ql.AddListener(sourceChan)
		}

	}
}

func (ql *QueueListener) AddListener(msgs <-chan amqp.Delivery) { //listen for messages from the channel
	for msg := range msgs {
		r := bytes.NewReader(msg.Body)
		d := gob.NewDecoder(r)
		sd := new(dto.InventoryMessage)
		d.Decode(sd)

		fmt.Printf("Received message: %v\n", sd)

		ed := EventData{
			Item:  sd.Item,
			Count: sd.Count,
		}

		ql.ea.PublishEvent("MessageReceived_"+msg.RoutingKey, ed)
	}
}

func (ql *QueueListener) DiscoverInventories() {
	ql.ch.ExchangeDeclare(
		queueutils.InventoryDiscoveryExchange,
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)

	ql.ch.Publish(
		queueutils.InventoryDiscoveryExchange,
		"",
		false,
		false,
		amqp.Publishing{},
	)
}
