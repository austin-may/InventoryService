package main

import (
	"bytes"
	"encoding/gob"
	"log"
	queueutils "my-go-apps/Inventory/queue-utils"
	"my-go-apps/inventory/dto"

	"github.com/streadway/amqp"
)

var url = "amqp://guest:guest@localhost:5672"

var name = "inventory"

type Shipping struct {
	destination string
	inventory   map[string]int
}

func main() {
	inventory1 := make(map[string]int)
	inventory1["broccoli"] = 40
	inventory1["apples"] = 75
	inventory1["kiwi"] = 22
	inventory1["kale"] = 18

	inventory2 := make(map[string]int)
	inventory2["potatoes"] = 40
	inventory2["toothpaste"] = 60
	inventory2["doritos"] = 23
	inventory2["carrots"] = 7

	inventory3 := make(map[string]int)
	inventory3["bagels"] = 90
	inventory3["burgers"] = 50
	inventory3["apple juice"] = 44
	inventory3["poptarts"] = 55963

	sites := []string{"Brookhaven", "Ansley Mall", "GA Tech Campus"}

	for _, site := range sites {
		conn, ch := queueutils.GetChannel(url)
		defer conn.Close()
		defer ch.Close()
		dataQueue := queueutils.GetQueue(site, ch, false) // will create the queue for us, if it doesn't already exist

		//discoveryQueue := queueutils.GetQueue("", ch, true)
		ch.QueueBind(
			dataQueue.Name,
			"",
			"amq.direct",
			false,
			nil,
		)

		inventoryToShip := make(map[string]int)
		switch site {
		case "Brookhaven":
			inventoryToShip = inventory1
		case "Ansley Mall":
			inventoryToShip = inventory2
		case "GA Tech Campus":
			inventoryToShip = inventory3
		}

		shippingRoute := Shipping{
			destination: site,
			inventory:   inventoryToShip,
		}

		//go listenForDiscoverRequests(site, ch)

		buf := new(bytes.Buffer)   //allows for reading/writing encoded data in memory
		enc := gob.NewEncoder(buf) // this enables the encoding

		for key, value := range shippingRoute.inventory {
			reading := dto.InventoryMessage{
				Item:  key,
				Count: value,
				Site:  site,
			}

			log.Printf("Delivering %d of %s", value, key)

			buf.Reset() // removes any previous data and resets the buffer back to its intiial position
			enc = gob.NewEncoder(buf)
			enc.Encode(reading) // this does the actual encoding

			msg := amqp.Publishing{
				Body: buf.Bytes(),
			}

			ch.Publish(
				"amq.direct",
				"",
				false,
				false,
				msg)
		}
	}

}

// func publishQueueName(ch *amqp.Channel) {
// 	msg := amqp.Publishing{Body: []byte(name)}
// 	ch.Publish(
// 		"amq.direct",
// 		"",
// 		false,
// 		false,
// 		msg,
// 	)
// }

// func listenForDiscoverRequests(name string, ch *amqp.Channel) {
// 	msgs, _ := ch.Consume(
// 		name,
// 		"",
// 		true,
// 		false,
// 		false,
// 		false,
// 		nil)

// 	for range msgs {
// 		publishQueueName(ch)
// 	}

// }
