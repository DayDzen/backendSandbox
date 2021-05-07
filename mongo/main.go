package main

import (
	"context"

	bs "github.com/rekamarket/mongodb-storage-lib/services"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func main() {
	s, err := NewDBStorage("mongodb://localhost:27017", "reka", "orders_groups")
	if err != nil {
		panic(err)
	}
	// m := &TestModel{}
	// m.A = 1
	// m.B = "b"
	// m.C = 1.1
	// if _, err := s.InsertOne(context.Background(), m); err != nil {
	// 	panic(err)
	// }
	// update
	// m := &TestModel{}
	// m.SetHexID("608bb62d2fba7eada6867724")
	// m.A = 3
	// m.B = "pzdc"
	// if err := s.UpdateOne(context.Background(), m); err != nil {
	// 	panic(err)
	// }

	var updateFields bson.D
	var uf bson.D

	oid, err := primitive.ObjectIDFromHex("608a8c1f08a12413346a083b")
	if err != nil {
		panic(err)
	}

	um := UpdateOrdersGroupByAdmin{
		ClientID: "123123123",
		ClientContact: &UpdateClientContactByAdmin{
			FirstName:  "FirstName",
			LastName:   "LastName",
			Patronymic: "Patronymic",
			Phone:      "Phone",
			Email:      "Email",
		},
	}

	updateFields = append(updateFields, bson.E{"client_contact.first_name", um.ClientContact.FirstName})
	updateFields = append(updateFields, bson.E{"client_contact.last_name", um.ClientContact.LastName})
	updateFields = append(updateFields, bson.E{"client_contact.patronymic", um.ClientContact.Patronymic})
	updateFields = append(updateFields, bson.E{"client_contact.phone", um.ClientContact.Phone})
	updateFields = append(updateFields, bson.E{"client_contact.email", um.ClientContact.Email})

	for _, v := range updateFields {
		if (v.Value != 0) && (v.Value != "") {
			uf = append(uf, v)
		}
	}

	update := bson.D{bson.E{"$set", uf}}
	if _, err := s.GetCollection().UpdateOne(context.Background(), bson.M{"_id": oid}, update); err != nil {
		panic(err)
	}
}

type UpdateOrdersGroupByAdmin struct {
	ClientID      string                      `bson:"client_id,omitempty"`
	ClientContact *UpdateClientContactByAdmin `bson:"client_contact,omitempty"`
}

type UpdateClientContactByAdmin struct {
	FirstName  string `bson:"first_name,omitempty"`
	LastName   string `bson:"last_name,omitempty"`
	Patronymic string `bson:"patronymic,omitempty"`
	Phone      string `bson:"phone,omitempty"`
	Email      string `bson:"email,omitempty"`
}

func NewDBStorage(mongoURI, dbName, collectionName string) (*bs.BaseStorage, error) {
	baseStorage, err := bs.NewBaseStorage(context.Background(), mongoURI, dbName, collectionName)
	if err != nil {
		return nil, err
	}
	return baseStorage, nil
}
