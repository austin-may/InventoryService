package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"
	"my-go-apps/InventoryService/graph/generated"
	"my-go-apps/InventoryService/graph/model"
	inventorydata "my-go-apps/InventoryService/inventory-data"
)

func (r *mutationResolver) CreateInventory(ctx context.Context, input model.NewInventory) (*model.NewInventoryResponse, error) {
	response, err := inventorydata.AddInventory(input)
	if err != nil {
		return nil, err
	}
	return response, nil
}

func (r *mutationResolver) UpdateInventory(ctx context.Context, input model.InventoryToUpdate) (*model.Inventory, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *mutationResolver) DeleteInventory(ctx context.Context, inventoryID *int) (*model.Inventory, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Inventory(ctx context.Context) ([]*model.Inventory, error) {
	inventory, err := inventorydata.GetInventory()
	if err != nil {
		return nil, err
	}
	return inventory, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
