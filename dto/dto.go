package dto

import "encoding/gob"

type InventoryMessage struct {
	Item  string
	Count int
	Site string
}

func init() {
	gob.Register(InventoryMessage{})
}
