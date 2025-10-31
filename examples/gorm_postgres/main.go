package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mhshajib/chronz/chronz"
	chronzgorm "github.com/mhshajib/chronz/chronz/gorm"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Order struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `tz:"local" gorm:"serializer:tztime"`
}

func main() {
	ctx := chronz.WithTZName(context.Background(), "Asia/Dhaka")

	dsn := "host=localhost user=root password=secret dbname=orders port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	schema.RegisterSerializer("tztime", chronzgorm.TZTimeSerializer{})

	_ = db.AutoMigrate(&Order{})

	_ = db.WithContext(ctx).Create(&Order{CreatedAt: time.Now()}).Error

	var out []Order
	_ = db.WithContext(ctx).Find(&out).Error

	fmt.Println("Orders (localized):")
	for _, o := range out {
		fmt.Println(" ->", o.CreatedAt)
	}
}
