package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func DBHandler(db *gorm.DB) handler {
	return handler{db}
}

func InitDB(dbURL string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	checkErr(err)

	db.AutoMigrate(&Job{})
	db.AutoMigrate(&Image{})

	return db
}
