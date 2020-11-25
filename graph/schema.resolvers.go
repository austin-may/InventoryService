package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"my-go-apps/InventoryService/graph/generated"
	"my-go-apps/InventoryService/graph/model"
	inventorydata "my-go-apps/InventoryService/inventory-data"
)

func (r *queryResolver) Inventory(ctx context.Context) ([]*model.Inventory, error) {
	inventory, err := inventorydata.GetInventory()
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
