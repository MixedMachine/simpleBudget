package repository

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

const (
	DB                    = "budgetdb"
	INCOME_COLLECTION     = "income"
	EXPENSE_COLLECTION    = "expense"
	ALLOCATION_COLLECTION = "allocation"
)

func CreateCollections(ctx *context.Context, client *mongo.Client) map[string]*Collection {
	incomeCollection := client.Database(DB).Collection(INCOME_COLLECTION)
	expenseCollection := client.Database(DB).Collection(EXPENSE_COLLECTION)
	allocationCollection := client.Database(DB).Collection(ALLOCATION_COLLECTION)

	ic := NewCollection(incomeCollection, *ctx)
	ec := NewCollection(expenseCollection, *ctx)
	ac := NewCollection(allocationCollection, *ctx)

	return map[string]*Collection{
		"income":     ic,
		"expense":    ec,
		"allocation": ac,
	}
}
