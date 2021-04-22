package inventorydata

import (
	"context"
	"fmt"
	datamanager "my-go-apps/InventoryService/data-manager"
	"my-go-apps/InventoryService/graph/model"
	"strconv"
	"time"
	"strings"
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

func GetNutritionFactsFromConsumedInventory(inventoryConsumed []*model.InventoryConsumed) ([]*model.NutritionFacts, error) {
	var nutritionNames []string
	nutritionNameAmountMap := make(map[string]int)
	for _, inventory := range inventoryConsumed {
		nutritionNames = append(nutritionNames, "'"+inventory.Name+"'")
		nutritionNameAmountMap[strings.ToLower(inventory.Name)] = inventory.Amount
	}
	nutritionString := strings.Join(nutritionNames, ",")
	nutritionQuery := fmt.Sprintf(`select i.Name, v.VitaminType, iv.PercentDailyValue
	from VitaminDB..Inventory i
	join VitaminDB..InventoryVitamin iv
	on i.InventoryID = iv.InventoryID
	join VitaminDB..Vitamin v
	on v.VitaminID = iv.VitaminID
	where i.Name in (%s)`, nutritionString)

	fmt.Printf(nutritionQuery)

	context, cancel := context.WithTimeout(context.Background(), 8000*time.Millisecond)
	defer cancel()

	nutritionResults, err := datamanager.DbConn.QueryContext(context, nutritionQuery)

	if (err != nil) {
		return nil, err
	}

	defer nutritionResults.Close()

	nutritionFactsList := make([]*model.NutritionFacts, 0)
	for nutritionResults.Next() {
		var nutritionFacts model.NutritionFacts
		var nutritionFact model.NutritionFact


		nutritionResults.Scan(&nutritionFacts.InventoryName, &nutritionFact.Vitamin, &nutritionFact.Percent)
		percentMultipliedByAmountConsumed := nutritionNameAmountMap[nutritionFacts.InventoryName] * nutritionFact.Percent //Amount * Percent
		nutritionFact.Percent = percentMultipliedByAmountConsumed
		nutritionFacts.NutritionFact = &nutritionFact
		fmt.Printf("%s", nutritionFact.Vitamin)

		nutritionFactsList = append(nutritionFactsList, &nutritionFacts)
	}
	return nutritionFactsList, nil
}

func GetNutritionInfoByInventoryId(inventoryId int) ([]*model.NutritionFacts, error) {
	fmt.Printf("the invetoryID is: %d", inventoryId)

	nutritionInfoQuery := fmt.Sprintf(`select i.Name, v.VitaminType, iv.PercentDailyValue
	from VitaminDB..InventoryVitamin iv
	join VitaminDB..Vitamin v 
	on v.VitaminID = iv.VitaminID
	join VitaminDB..Inventory i
	on i.InventoryID = iv.InventoryID
	where i.InventoryID = %d`, inventoryId)

	context, cancel := context.WithTimeout(context.Background(), time.Millisecond * 8000)
	defer cancel()

	nutritionInfoResults, err := datamanager.DbConn.QueryContext(context, nutritionInfoQuery)

	if (err != nil) {
		return nil, err
	}

	defer nutritionInfoResults.Close()

	
	nutritionFactsList := make([]*model.NutritionFacts, 0)
	for nutritionInfoResults.Next() {
		var nutritionFacts model.NutritionFacts
		var nutritionFact model.NutritionFact
		nutritionInfoResults.Scan(&nutritionFacts.InventoryName, &nutritionFact.Vitamin, &nutritionFact.Percent)
		nutritionFacts.NutritionFact = &nutritionFact

		nutritionFactsList = append(nutritionFactsList, &nutritionFacts)
	
		fmt.Printf("%s", nutritionFacts)
	}
	
	return nutritionFactsList, nil
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
