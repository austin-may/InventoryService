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
	inventory2["potato"] = 40
	inventory2["lemon"] = 60
	inventory2["orange"] = 23
	inventory2["carrot"] = 7

	sites := []string{"Brookhaven", "Ansley Mall", "GA Tech Campus"}

	//region Brookhaven Inventory
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

	//endregion

	//region Ansley Mall Inventory
	potatoNutritionSlice := make([]map[string]int, 0)
	potatoNutritionSlice = append(potatoNutritionSlice, map[string]int{"Potassium": 25})
	potatoNutritionSlice = append(potatoNutritionSlice, map[string]int{"Sodium": 1})
	potatoNutritionSlice = append(potatoNutritionSlice, map[string]int{"A": 0})
	potatoNutritionSlice = append(potatoNutritionSlice, map[string]int{"C": 32})
	potatoNutritionSlice = append(potatoNutritionSlice, map[string]int{"Calcium": 2})
	potatoNutritionSlice = append(potatoNutritionSlice, map[string]int{"Iron": 6})
	potatoNutritionSlice = append(potatoNutritionSlice, map[string]int{"B6": 20})
	potatoNutritionSlice = append(potatoNutritionSlice, map[string]int{"Magnesium": 11})
	potatoNutritionSlice = append(potatoNutritionSlice, map[string]int{"D": 0})
	potatoNutritionSlice = append(potatoNutritionSlice, map[string]int{"B12": 0})

	potato := dto.Item{
		Name:           "potato",
		NutritionFacts: potatoNutritionSlice,
	}

	lemonNutritionSlice := make([]map[string]int, 0)
	lemonNutritionSlice = append(lemonNutritionSlice, map[string]int{"Potassium": 2})
	lemonNutritionSlice = append(lemonNutritionSlice, map[string]int{"Sodium": 0})
	lemonNutritionSlice = append(lemonNutritionSlice, map[string]int{"A": 0})
	lemonNutritionSlice = append(lemonNutritionSlice, map[string]int{"C": 51})
	lemonNutritionSlice = append(lemonNutritionSlice, map[string]int{"Calcium": 2})
	lemonNutritionSlice = append(lemonNutritionSlice, map[string]int{"Iron": 2})
	lemonNutritionSlice = append(lemonNutritionSlice, map[string]int{"B6": 0})
	lemonNutritionSlice = append(lemonNutritionSlice, map[string]int{"Magnesium": 1})
	lemonNutritionSlice = append(lemonNutritionSlice, map[string]int{"D": 0})
	lemonNutritionSlice = append(lemonNutritionSlice, map[string]int{"B12": 0})

	lemon := dto.Item{
		Name:           "lemon",
		NutritionFacts: lemonNutritionSlice,
	}

	orangeNutritionSlice := make([]map[string]int, 0)
	orangeNutritionSlice = append(orangeNutritionSlice, map[string]int{"Potassium": 7})
	orangeNutritionSlice = append(orangeNutritionSlice, map[string]int{"Sodium": 0})
	orangeNutritionSlice = append(orangeNutritionSlice, map[string]int{"A": 6})
	orangeNutritionSlice = append(orangeNutritionSlice, map[string]int{"C": 116})
	orangeNutritionSlice = append(orangeNutritionSlice, map[string]int{"Calcium": 5})
	orangeNutritionSlice = append(orangeNutritionSlice, map[string]int{"Iron": 1})
	orangeNutritionSlice = append(orangeNutritionSlice, map[string]int{"B6": 0})
	orangeNutritionSlice = append(orangeNutritionSlice, map[string]int{"Magnesium": 3})
	orangeNutritionSlice = append(orangeNutritionSlice, map[string]int{"D": 0})
	orangeNutritionSlice = append(orangeNutritionSlice, map[string]int{"B12": 0})

	orange := dto.Item{
		Name:           "orange",
		NutritionFacts: orangeNutritionSlice,
	}

	carrotNutritionSlice := make([]map[string]int, 0)
	orangeNutritionSlice = append(orangeNutritionSlice, map[string]int{"Potassium": 7})
	carrotNutritionSlice = append(carrotNutritionSlice, map[string]int{"A": 93})
	carrotNutritionSlice = append(carrotNutritionSlice, map[string]int{"C": 7})
	carrotNutritionSlice = append(carrotNutritionSlice, map[string]int{"K": 11})
	carrotNutritionSlice = append(carrotNutritionSlice, map[string]int{"Iron": 2})
	carrotNutritionSlice = append(carrotNutritionSlice, map[string]int{"B6": 8})
	carrotNutritionSlice = append(carrotNutritionSlice, map[string]int{"Magnesium": 3})

	carrot := dto.Item{
		Name:           "carrot",
		NutritionFacts: carrotNutritionSlice,
	}

	inventoryItemsForSite = append(inventoryItemsForSite, potato)
	inventoryItemsForSite = append(inventoryItemsForSite, lemon)
	inventoryItemsForSite = append(inventoryItemsForSite, orange)
	inventoryItemsForSite = append(inventoryItemsForSite, carrot)
	//endregion

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
