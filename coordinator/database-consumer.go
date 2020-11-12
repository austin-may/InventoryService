package coordinator

import (
	"bytes"
	"encoding/gob"
	"my-go-apps/Inventory/dto"
	queueutils "my-go-apps/Inventory/queue-utils"
	"time"

	"github.com/streadway/amqp"
)

/*
This package is for consuming events that may be *of interest* to certain consumers
This is an intermediate/advanced approached to my distributed design that I will revist later
*/

type DatabaseConsumer struct {
	er      EventRaiser
	conn    *amqp.Connection
	ch      *amqp.Channel
	queue   *amqp.Queue
	sources []string
}

const maxRate = 5 * time.Second

func NewDatabaseConsumer(er EventRaiser) *DatabaseConsumer {
	dc := DatabaseConsumer{
		er: er,
	}

	dc.conn, dc.ch = queueutils.GetChannel(url)
	dc.queue = queueutils.GetQueue(queueutils.PersistReadingsQueue, dc.ch, false)

	dc.er.AddListener("DataSourceDiscovered", func(eventData interface{}) {
		dc.SubscribeToDataEvent(eventData.(string))
	})

	return &dc
}

func (dc *DatabaseConsumer) SubscribeToDataEvent(eventName string) {
	for _, v := range dc.sources {
		if v == eventName {
			return
		}
	}

	//From the transcript:
	// For the callback, I'm going to pass in a self-executing function that will return the callback itself.
	// The reason for this will be apparent in a second.
	// To make the function self-executing, I need to invoke it by adding the mashed parens after the definition.
	// What this is going to buy me is a new isolated variable scope that is going be created every time I call this function.
	// This is a trick that is going to allow the event handlers to register a bit of state that
	// I'm going to need in order to throttle down the rate that the messages are coming in from the event sources.
	dc.er.AddListener("MessageReceived_"+eventName, func() func(interface{}) {
		prevTime := time.Unix(0, 0)

		buf := new(bytes.Buffer)

		return func(eventData interface{}) {
			ed := eventData.(EventData)
			if time.Since(prevTime) > maxRate {
				prevTime = time.Now()

				sm := dto.InventoryMessage{
					Item:  ed.Item,
					Count: ed.Count,
				}

				buf.Reset()

				enc := gob.NewEncoder(buf)
				enc.Encode(sm)

				msg := amqp.Publishing{
					Body: buf.Bytes(),
				}

				dc.ch.Publish(
					"",
					queueutils.PersistReadingsQueue,
					false,
					false,
					msg)

			}
		}
	}())
}
