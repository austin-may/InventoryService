package datamanager

import (
	"fmt"
	"my-go-apps/InventoryService/dto"
	"strings"
)

var inventory map[string]int

func InsertInventoryItem(reading *dto.InventoryMessage) error {

	var lastInsertId int
	insertInventory :=
		fmt.Sprintf("INSERT INTO Inventory (Name, Count, Site) "+
			"OUTPUT INSERTED.[InventoryID] "+
			"VALUES ('%s', %d, '%s')", reading.Item.Name, reading.Count, reading.Site)

	err := DbConn.QueryRow(insertInventory).Scan(&lastInsertId)

	if err != nil {
		println("Error:", err.Error())
	} else {
		println("LastInsertId:", lastInsertId)
	}

	//loop through nutrition facts
	for _, val := range reading.Item.NutritionFacts {
		var vitaminTypes []string
		for vitaminType, _ := range val {
			vitaminTypes = append(vitaminTypes, fmt.Sprintf("'%s'", vitaminType))
		}
		q :=
			fmt.Sprintf("SELECT VitaminID, VitaminType "+
				"FROM VitaminDB..Vitamin "+
				"WHERE VitaminType in (%s)", strings.Join(vitaminTypes, ","))

		rows, err := DbConn.Query(q)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var vitaminId int
			var vitaminType string

			rows.Scan(&vitaminId, &vitaminType)

			insertInventoryVitamin :=
				fmt.Sprintf("INSERT INTO InventoryVitamin (InventoryID, VitaminID, PercentDailyValue) "+
					"VALUES (%d, %d, %d)", lastInsertId, vitaminId, val[vitaminType])

			_, err = DbConn.Exec(insertInventoryVitamin)
		}
	}

	return err
}

func UpdateInventoryItem(reading *dto.InventoryMessage) error {

	q :=
		fmt.Sprintf("UPDATE Inventory "+
			"SET Count = %d "+
			"WHERE Name = '%s' and Site = '%s'", reading.Count, reading.Item.Name, reading.Site)

	_, err := DbConn.Exec(q)

	return err
}

func GetExistingInventory() (map[string]int, error) {
	inventory = make(map[string]int)
	q := `
		SELECT Name, Count
		FROM VitaminDB..Inventory
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
