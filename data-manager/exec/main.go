package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	datamanager "my-go-apps/Inventory/data-manager"
	"my-go-apps/Inventory/dto"
	queueutils "my-go-apps/Inventory/queue-utils"
)

const url = "amqp://guest:guest@localhost:5672"

func main() {
	fmt.Println("Started listening for messages to arrive on queues...")
	sites := []string{"Brookhaven", "Ansley Mall", "GA Tech Campus"}

	for _, site := range sites {
		conn, ch := queueutils.GetChannel(url)
		defer conn.Close()
		defer ch.Close()

		msgs, err := ch.Consume(
			site,
			"",
			false,
			true,
			false,
			false,
			nil)

		if err != nil {
			log.Fatalln("Failed to get access to messages")
		}

		//putting this here requires restarting the listener to make the incoming items we've already stored idempotent
		//we'll look to correct that later
		inventory, err := datamanager.GetExistingInventory()
		if err == nil {
			for msg := range msgs {
				buf := bytes.NewReader(msg.Body)
				dec := gob.NewDecoder(buf)
				sd := &dto.InventoryMessage{}
				dec.Decode(sd)

				fmt.Printf("Received message: %v\n", sd)

				//determine if it's existing inventory or not
				if _, ok := inventory[sd.Item]; ok {
					err = datamanager.UpdateInventoryItem(sd)
				} else {
					err = datamanager.InsertInventoryItem(sd)
				}
				if err != nil {
					log.Printf("Failed to save reading from inventory %v. Error: %s", sd.Item, err.Error())
				} else {
					msg.Ack(true)
				}
			}
		}
	}
}
