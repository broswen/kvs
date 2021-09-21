package db

import (
	"errors"
	"fmt"
	"os"

	"github.com/broswen/kvs/internal/item"
	items "github.com/broswen/kvs/internal/item"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBService struct {
	db *gorm.DB
}

func New() (DBService, error) {
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	user := os.Getenv("POSTGRES_USER")
	pass := os.Getenv("POSTGRES_PASS")
	dbname := os.Getenv("POSTGRES_DB")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s  sslmode=disable", host, user, pass, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return DBService{}, err
	}
	db.AutoMigrate(&item.Item{})
	return DBService{
		db: db,
	}, nil
}

func (dbs DBService) Get(key string) (item.Item, error) {
	item := items.Item{Key: key}
	tx := dbs.db.First(&item)
	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return items.Item{}, items.ErrItemNotFound{Err: tx.Error}
		}
		return items.Item{}, items.ErrService{Err: tx.Error}
	}
	return item, nil
}

func (dbs DBService) Set(item item.Item) error {
	tx := dbs.db.Save(item)
	if tx.Error != nil {
		return items.ErrService{Err: tx.Error}
	}
	return nil
}
