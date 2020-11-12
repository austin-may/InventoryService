package datamanager

import (
	"fmt"
	"my-go-apps/Inventory/dto"
)

var inventory map[string]int

func InsertInventoryItem(reading *dto.InventoryMessage) error {

	q :=
		fmt.Sprintf("INSERT INTO InventoryReserve (Name, Count, Site) "+
			"VALUES ('%s', %d, '%s')", reading.Item, reading.Count, reading.Site)

	_, err := DbConn.Exec(q)

	return err
}

func UpdateInventoryItem(reading *dto.InventoryMessage) error {

	q :=
		fmt.Sprintf("UPDATE InventoryReserve "+
			"SET Count = %d"+
			"WHERE Name = '%s' and Site = '%s'", reading.Count, reading.Item, reading.Site)

	_, err := DbConn.Exec(q)

	return err
}

func GetExistingInventory() (map[string]int, error) {
	inventory = make(map[string]int)
	q := `
		SELECT Name, Count
		FROM InventoryDB..InventoryReserve
	`

	rows, err := DbConn.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var count int

		rows.Scan(&name, &count)

		inventory[name] = count
	}

	return inventory, nil
}
