package db

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/broswen/kvs/internal/item"
	items "github.com/broswen/kvs/internal/item"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type DBService struct {
	db *gorm.DB
}

func New() (*DBService, error) {
	sqliteDB, ok := os.LookupEnv("SQLITE_DB")
	var db *gorm.DB
	var err error
	if ok {
		// sqlite db path is set, use sqlite driver
		log.Printf("using sqlite: %v\n", sqliteDB)
		db, err = gorm.Open(sqlite.Open(sqliteDB), &gorm.Config{})
	} else {
		postgresHost, ok := os.LookupEnv("POSTGRES_HOST")
		if !ok {
			return &DBService{}, errors.New("either SQLITE_DB or POSTGRES_HOST must be set")
		}
		postgresPort := os.Getenv("POSTGRES_PORT")
		postgresUser := os.Getenv("POSTGRES_USER")
		postgresPass := os.Getenv("POSTGRES_PASS")
		postgresDBName := os.Getenv("POSTGRES_DB")
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s  sslmode=disable", postgresHost, postgresUser, postgresPass, postgresDBName, postgresPort)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		return &DBService{}, err
	}
	db.AutoMigrate(&item.Item{})
	return &DBService{
		db: db,
	}, nil
}

func (dbs *DBService) IsNil() bool {
	return dbs == nil
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

func (dbs DBService) Delete(key string) error {
	item := items.Item{Key: key}
	tx := dbs.db.Delete(&item)
	if tx.Error != nil {
		return items.ErrService{Err: tx.Error}
	}
	return nil
}
