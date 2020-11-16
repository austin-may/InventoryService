package datamanager

import (
	"fmt"
	"my-go-apps/InventoryService/dto"
)

var inventory map[string]int

func InsertInventoryItem(reading *dto.InventoryMessage) error {

	var lastInsertId int
	insertInventory :=
		fmt.Sprintf("INSERT INTO Inventory (Name, Count, Site) "+
			"OUTPUT INSERTED.[ItemID] "+
			"VALUES ('%s', %d, '%s')", reading.Item.Name, reading.Count, reading.Site)

	err := DbConn.QueryRow(insertInventory).Scan(&lastInsertId)

	if err != nil {
		println("Error:", err.Error())
	} else {
		println("LastInsertId:", lastInsertId)
	}

	//loop through nutrition facts
	for _, val := range reading.Item.NutritionFacts {
		//this is super chatty, need to find a way to do this once
		for vitaminType, percentDailyValue := range val {
			q :=
				fmt.Sprintf("SELECT VitaminID "+
					"FROM VitaminDB..Vitamin "+
					"WHERE VitaminType = '%s'", vitaminType)

			rows, err := DbConn.Query(q)
			if err != nil {
				return err
			}
			defer rows.Close()

			for rows.Next() {
				var vitaminId int

				rows.Scan(&vitaminId)

				insertInventoryVitamin :=
					fmt.Sprintf("INSERT INTO InventoryVitamin (ItemID, VitaminID, PercentDailyValue) "+
						"VALUES (%d, %d, %d)", lastInsertId, vitaminId, percentDailyValue)

				_, err = DbConn.Exec(insertInventoryVitamin)
			}
		}
	}

	return err
}

func UpdateInventoryItem(reading *dto.InventoryMessage) error {

	q :=
		fmt.Sprintf("UPDATE Inventory "+
			"SET Count = %d"+
			"WHERE Name = '%s' and Site = '%s'", reading.Count, reading.Item, reading.Site)

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
