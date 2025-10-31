package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mhshajib/chronz"
	chronzgorm "github.com/mhshajib/chronz/chronz_gorm"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type Order struct {
	ID        uint      `gorm:"primaryKey"`
	CreatedAt time.Time `tz:"local" gorm:"serializer:tztime"`
}

func main() {
	// Choose the request timezone
	ctx := chronz.WithTZName(context.Background(), "Asia/Dhaka")

	// DB connect
	dsn := "host=localhost user=root password=secret dbname=orders port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	// Register the tz serializer once at boot
	schema.RegisterSerializer("tztime", chronzgorm.TZTimeSerializer{})

	// Migrate + insert (local -> UTC)
	_ = db.AutoMigrate(&Order{})
	_ = db.WithContext(ctx).Create(&Order{CreatedAt: time.Now()}).Error

	// Query (UTC -> local)
	var out []Order
	_ = db.WithContext(ctx).Find(&out).Error

	// Example raw WHERE using ArgTime (local -> UTC)
	_ = db.WithContext(ctx).
		Where("created_at >= @created_at", chronzgorm.ArgTime(ctx, "created_at", time.Now().Add(-24*time.Hour))).
		Find(&out).Error

	fmt.Println("Orders (localized):")
	for _, o := range out {
		fmt.Println(" ->", o.CreatedAt)
	}
}
