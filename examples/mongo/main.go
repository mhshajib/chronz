package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mhshajib/chronz/chronz"
	chronzmongo "github.com/mhshajib/chronz/chronz/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Order struct {
	ID        string    `bson:"_id,omitempty"`
	CreatedAt time.Time `bson:"created_at" tz:"local"`
}

func main() {
	ctx := chronz.WithTZName(context.Background(), "Asia/Dhaka")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	base := client.Database("orders_db").Collection("orders")
	coll := chronzmongo.WrapCollection(base)

	_, _ = coll.InsertOne(ctx, Order{CreatedAt: time.Now()})
	res := coll.FindOne(ctx, bson.M{})
	var out Order
	_ = chronzmongo.DecodeLocal(ctx, res, &out)

	fmt.Println("Inserted & read (localized):", out.CreatedAt)
}
