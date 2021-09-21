package item

import (
	"errors"
	"log"

	"honnef.co/go/tools/config"
)

type Item struct {
	Key   string `gorm:"primaryKey"`
	Value string
}

type ErrItemNotFound struct {
	Err error
}

func (e ErrItemNotFound) Unwrap() error {
	return e.Err
}

func (e ErrItemNotFound) Error() string {
	return e.Err.Error()
}

type ErrService struct {
	Err error
}

func (e ErrService) Unwrap() error {
	return e.Err
}

func (e ErrService) Error() string {
	return e.Err.Error()
}

type ItemGetterSetter interface {
	Get(string) (Item, error)
	Set(Item) error
}

type ItemService struct {
	cache  ItemGetterSetter
	db     ItemGetterSetter
	config config.Config
}

func New(cache ItemGetterSetter, db ItemGetterSetter) (ItemService, error) {

	return ItemService{
		cache: cache,
		db:    db,
	}, nil
}

func (is ItemService) Get(key string) (Item, error) {
	item, err := is.cache.Get(key)
	if err == nil {
		// found in cache, return
		return item, nil
	}

	if errors.As(err, &ErrService{}) {
		log.Printf("get cache: %v\n", err)
	}

	if errors.As(err, &ErrItemNotFound{}) {
		log.Printf("cache miss: %v\n", err)
	}

	item, err = is.db.Get(key)
	if err != nil {
		// couldn't get from db, either err or not found
		return Item{}, err
	}
	// item was found in db, set in cache
	err = is.cache.Set(item)
	if err != nil {
		log.Printf("set cache: %v\n", err)
	}

	// return found item
	return item, nil
}

func (is ItemService) Set(item Item) error {
	// set in db
	err := is.db.Set(item)
	if err != nil {
		return ErrService{err}
	}
	// set in cache
	err = is.cache.Set(item)
	if err != nil {
		return ErrService{err}
	}

	return nil
}
