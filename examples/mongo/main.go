package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mhshajib/chronz"
	chronzMongo "github.com/mhshajib/chronz/chronz_mongo"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Order struct {
	ID        string    `bson:"_id,omitempty"`
	CreatedAt time.Time `bson:"created_at" tz:"local"`
}

func main() {
	// Choose the request timezone
	ctx := chronz.WithTZName(context.Background(), "Asia/Dhaka")

	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(ctx)

	base := client.Database("orders_db").Collection("orders")
	coll := chronzMongo.WrapCollection(base)

	// Insert (local -> UTC)
	_, _ = coll.InsertOne(ctx, Order{CreatedAt: time.Now()})

	// FindOne (UTC -> local)
	res := coll.FindOne(ctx, bson.M{})
	var out Order
	_ = chronzMongo.DecodeLocal(ctx, res, &out)

	fmt.Println("Inserted & read (localized):", out.CreatedAt)

	// Example aggregate with NormalizePipeline
	pipe := chronzMongo.NormalizePipeline(ctx, mongo.Pipeline{
		{{"$match", bson.M{"created_at": bson.M{"$gte": time.Now().Add(-24 * time.Hour)}}}},
	})
	cur, _ := coll.Aggregate(ctx, pipe)
	_ = cur.All(ctx, &[]Order{})
}
