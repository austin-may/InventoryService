package main

import (
	"fmt"
	"my-go-apps/InventoryService/coordinator"
)

var dc *coordinator.DatabaseConsumer //package level variable

func main() {
	ea := coordinator.NewEventAggregator()
	dc = coordinator.NewDatabaseConsumer(ea)
	ql := coordinator.NewQueueListener(ea)
	go ql.ListenForNewSource()

	var a string
	fmt.Scanln(&a)
}
