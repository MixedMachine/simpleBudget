package repository

import (
	. "github.com/mixedmachine/simple-budget-app/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"context"
)

var (
	MONGO_URI string
)

type ContactCollection struct {
	client *mongo.Collection
	ctx    context.Context
}

func NewContactCollection(contactcollection *mongo.Collection, ctx context.Context) *ContactCollection {
	return &ContactCollection{
		client: contactcollection,
		ctx:    ctx,
	}
}

type Collection struct {
	client *mongo.Collection
	ctx    context.Context
}

func NewCollection(collection *mongo.Collection, ctx context.Context) *Collection {
	return &Collection{
		client: collection,
		ctx:    ctx,
	}
}

func (cc *ContactCollection) CreateContact(contact *Contact) error {
	_, err := cc.client.InsertOne(cc.ctx, contact)
	return err

}

func (c *ContactCollection) GetAllContacts() ([]Contact, error) {
	var contacts []Contact
	cursor, err := c.client.Find(c.ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(c.ctx, &contacts); err != nil {
		return nil, err
	}
	return contacts, nil
}

func (c *ContactCollection) UpdateContact(contact Contact) error {
	filter := bson.M{"_id": contact.ID}
	_, err := c.client.ReplaceOne(c.ctx, filter, contact)
	return err
}

func (c *ContactCollection) DeleteContact(id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := c.client.DeleteOne(c.ctx, filter)
	return err
}

func Create[T any](c *Collection, element T) error {
	_, err := c.client.InsertOne(c.ctx, element)
	return err

}

func GetAll[T any](c *Collection, elements T) error {
	cursor, err := c.client.Find(c.ctx, bson.M{})
	if err != nil {
		return err
	}
	if err = cursor.All(c.ctx, elements); err != nil {
		return err
	}
	return nil
}

func Get[T any](c *Collection, id primitive.ObjectID, element T) error {
	filter := bson.M{"_id": id}
	err := c.client.FindOne(c.ctx, filter).Decode(element)
	return err
}

func Update[T any](c *Collection, id primitive.ObjectID, element T) error {
	filter := bson.M{"_id": id}
	_, err := c.client.ReplaceOne(c.ctx, filter, element)
	return err
}

func Delete(c *Collection, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	_, err := c.client.DeleteOne(c.ctx, filter)
	return err
}
