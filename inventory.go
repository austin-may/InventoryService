package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"my-go-apps/InventoryService/dto"
	queueutils "my-go-apps/InventoryService/queue-utils"

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
	inventory1["broccoli"] = 99
	inventory1["apple"] = 75
	inventory1["kiwi"] = 105
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

	sites := []string{"Brookhaven" /*, "Ansley Mall", "GA Tech Campus"*/}

	broccoliNutritionSlice := make([]map[string]int, 0)
	broccoliNutritionSlice = append(broccoliNutritionSlice, map[string]int{"Potassium": 13})
	broccoliNutritionSlice = append(broccoliNutritionSlice, map[string]int{"Sodium": 2})
	broccoliNutritionSlice = append(broccoliNutritionSlice, map[string]int{"A": 18})
	broccoliNutritionSlice = append(broccoliNutritionSlice, map[string]int{"C": 220})
	broccoliNutritionSlice = append(broccoliNutritionSlice, map[string]int{"Calcium": 7})
	broccoliNutritionSlice = append(broccoliNutritionSlice, map[string]int{"Iron": 6})
	broccoliNutritionSlice = append(broccoliNutritionSlice, map[string]int{"B6": 15})
	broccoliNutritionSlice = append(broccoliNutritionSlice, map[string]int{"Magnesium": 7})
	broccoliNutritionSlice = append(broccoliNutritionSlice, map[string]int{"D": 0})
	broccoliNutritionSlice = append(broccoliNutritionSlice, map[string]int{"Cobalamin": 0})

	broccoli := dto.Item{
		Name:           "broccoli",
		NutritionFacts: broccoliNutritionSlice,
	}

	appleNutritionSlice := make([]map[string]int, 0)
	appleNutritionSlice = append(appleNutritionSlice, map[string]int{"Potassium": 5})
	appleNutritionSlice = append(appleNutritionSlice, map[string]int{"Sodium": 0})
	appleNutritionSlice = append(appleNutritionSlice, map[string]int{"A": 1})
	appleNutritionSlice = append(appleNutritionSlice, map[string]int{"C": 14})
	appleNutritionSlice = append(appleNutritionSlice, map[string]int{"Calcium": 1})
	appleNutritionSlice = append(appleNutritionSlice, map[string]int{"Iron": 1})
	appleNutritionSlice = append(appleNutritionSlice, map[string]int{"B6": 5})
	appleNutritionSlice = append(appleNutritionSlice, map[string]int{"Magnesium": 2})
	appleNutritionSlice = append(appleNutritionSlice, map[string]int{"D": 0})
	appleNutritionSlice = append(appleNutritionSlice, map[string]int{"Cobalamin": 0})

	apple := dto.Item{
		Name:           "apple",
		NutritionFacts: appleNutritionSlice,
	}

	kiwiNutritionSlice := make([]map[string]int, 0)
	kiwiNutritionSlice = append(kiwiNutritionSlice, map[string]int{"Potassium": 6})
	kiwiNutritionSlice = append(kiwiNutritionSlice, map[string]int{"Sodium": 0})
	kiwiNutritionSlice = append(kiwiNutritionSlice, map[string]int{"A": 1})
	kiwiNutritionSlice = append(kiwiNutritionSlice, map[string]int{"C": 106})
	kiwiNutritionSlice = append(kiwiNutritionSlice, map[string]int{"Calcium": 2})
	kiwiNutritionSlice = append(kiwiNutritionSlice, map[string]int{"Iron": 1})
	kiwiNutritionSlice = append(kiwiNutritionSlice, map[string]int{"B6": 0})
	kiwiNutritionSlice = append(kiwiNutritionSlice, map[string]int{"Magnesium": 3})
	kiwiNutritionSlice = append(kiwiNutritionSlice, map[string]int{"D": 0})
	kiwiNutritionSlice = append(kiwiNutritionSlice, map[string]int{"Cobalamin": 0})

	kiwi := dto.Item{
		Name:           "kiwi",
		NutritionFacts: kiwiNutritionSlice,
	}

	kaleNutritionSlice := make([]map[string]int, 0)
	kaleNutritionSlice = append(kaleNutritionSlice, map[string]int{"A": 133})
	kaleNutritionSlice = append(kaleNutritionSlice, map[string]int{"C": 134})
	kaleNutritionSlice = append(kaleNutritionSlice, map[string]int{"Calcium": 10})
	kaleNutritionSlice = append(kaleNutritionSlice, map[string]int{"Iron": 5})
	kaleNutritionSlice = append(kaleNutritionSlice, map[string]int{"B6": 10})
	kaleNutritionSlice = append(kaleNutritionSlice, map[string]int{"Magnesium": 7})

	kale := dto.Item{
		Name:           "kale",
		NutritionFacts: kaleNutritionSlice,
	}

	var inventoryItemsForSite = []dto.Item{}
	inventoryItemsForSite = append(inventoryItemsForSite, broccoli)
	inventoryItemsForSite = append(inventoryItemsForSite, apple)
	inventoryItemsForSite = append(inventoryItemsForSite, kiwi)
	inventoryItemsForSite = append(inventoryItemsForSite, kale)

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

		reading := []dto.InventoryMessage{}
		for key, value := range shippingRoute.inventory {
			for _, inventoryItem := range inventoryItemsForSite {
				if key == inventoryItem.Name {
					reading = append(reading, dto.InventoryMessage{
						Item:  inventoryItem,
						Site:  site,
						Count: value,
					})
				}
			}
		}
		for _, item := range reading {

			log.Printf("Delivering %d of %s", item.Count, item.Item.Name)

			buf.Reset() // removes any previous data and resets the buffer back to its intiial position
			enc = gob.NewEncoder(buf)
			enc.Encode(item) // this does the actual encoding

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
