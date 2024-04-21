package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Reciever struct {
	ID       int64
	Nickname string
}

type Message struct {
	ID   int
	Text string
}

var DB *gorm.DB

func ConnectDB() {

	dsn := "host=localhost user=postgres password=123 dbname=postgres port=5432 sslmode=disable TimeZone=Europe/Moscow"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	err = database.AutoMigrate(&Reciever{}, &Message{})

	if err != nil {
		return
	}

	DB = database
}
