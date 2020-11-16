package dto

import (
	"encoding/gob"
)

type Item struct {
	Name           string
	NutritionFacts []map[string]int
}

type InventoryMessage struct {
	Item  Item
	Count int
	Site  string
}

func init() {
	gob.Register(InventoryMessage{})
}
