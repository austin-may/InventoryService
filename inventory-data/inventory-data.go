package inventorydata

import (
	"context"
	"fmt"
	datamanager "my-go-apps/InventoryService/data-manager"
	"my-go-apps/InventoryService/graph/model"
	"strconv"
	"time"
)

func GetInventory() ([]*model.Inventory, error) {
	inventoryQuery := fmt.Sprintf("SELECT InventoryID, Name, Count, Price, ExpirationDate, Site, SkuNumber " +
		"FROM Inventory")

	inventoryVitaminQuery := fmt.Sprintf("SELECT InventoryVitaminID, InventoryID, VitaminID, PercentDailyValue " +
		"FROM InventoryVitamin")

	context, cancel := context.WithTimeout(context.Background(), 8000*time.Millisecond)
	defer cancel() //destroy resources associated with the timeout

	inventoryResults, err := datamanager.DbConn.QueryContext(context, inventoryQuery)
	if err != nil {
		return nil, err
	}

	inventoryVitaminResults, err := datamanager.DbConn.Query(inventoryVitaminQuery)
	if err != nil {
		return nil, err
	}

	defer inventoryResults.Close()
	defer inventoryVitaminResults.Close()

	inventoryList := make([]*model.Inventory, 0)
	for inventoryResults.Next() {
		var inventory model.Inventory
		inventoryResults.Scan(&inventory.InventoryID, &inventory.Name, &inventory.Count, &inventory.Price, &inventory.ExpirationDate, &inventory.Site, &inventory.SkuNumber)

		inventoryList = append(inventoryList, &inventory)
	}

	inventoryVitaminList := make([]*model.InventoryVitamin, 0)
	for inventoryVitaminResults.Next() {
		var inventoryVitamin model.InventoryVitamin
		inventoryVitaminResults.Scan(&inventoryVitamin.InventoryVitaminID, &inventoryVitamin.InventoryID, &inventoryVitamin.VitaminID, &inventoryVitamin.PercentDailyValue)

		inventoryVitaminList = append(inventoryVitaminList, &inventoryVitamin)
	}

	for _, inventoryElement := range inventoryList {
		for _, inventoryVitaminElement := range inventoryVitaminList {
			inventoryId, _ := strconv.Atoi(inventoryElement.InventoryID)
			if inventoryId == inventoryVitaminElement.InventoryID {
				inventoryElement.InventoryVitamin = append(inventoryElement.InventoryVitamin, inventoryVitaminElement)
			}
		}
	}

	return inventoryList, nil
}

func AddInventory(inventory model.NewInventory) (*model.NewInventoryResponse, error) {
	command := fmt.Sprintf(`INSERT INTO Inventory (Name, Count, Price, ExpirationDate, Site, SkuNumber) VALUES ('%s', %d, %f, '%s', '%s', '%s');
	SELECT InventoryID, Name, Count, Price, ExpirationDate, Site, SkuNumber FROM Inventory WHERE InventoryID = SCOPE_IDENTITY()`, inventory.Name, inventory.Count, inventory.Price, inventory.ExpirationDate, inventory.Site, inventory.SkuNumber)
	var addedItem model.NewInventoryResponse
	err := datamanager.DbConn.QueryRow(command).Scan(&addedItem.ID, &addedItem.Name, &addedItem.Count, &addedItem.Price, &addedItem.ExpirationDate, &addedItem.Site, &addedItem.SkuNumber)
	fmt.Println(command)
	if err != nil {
		return nil, err
	}
	fmt.Println(addedItem)
	return &addedItem, nil
}
