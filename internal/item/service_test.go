package item

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
)

type MapItemGetterSetter struct {
	data map[string]string
}

func newMapItemGetterSetter() *MapItemGetterSetter {
	return &MapItemGetterSetter{
		data: make(map[string]string),
	}
}

func (m MapItemGetterSetter) Get(key string) (Item, error) {
	val, ok := m.data[key]
	if !ok {
		return Item{}, errors.New("item not found")
	}
	return Item{Key: key, Value: val}, nil
}

func (m MapItemGetterSetter) Set(item Item) error {
	m.data[item.Key] = item.Value
	return nil
}

func (m MapItemGetterSetter) Delete(key string) error {
	delete(m.data, key)
	return nil
}

func (m *MapItemGetterSetter) IsNil() bool {
	return m == nil
}

func TestErrors(t *testing.T) {
	err := ErrService{errors.New("test")}
	if !errors.As(err, &ErrService{}) {
		t.Fatalf("error should be instance of ErrService")
	}
	err1 := ErrItemNotFound{errors.New("test")}
	if !errors.As(err1, &ErrItemNotFound{}) {
		t.Fatalf("error should be instance of ErrItemNotFound")
	}
}

func TestNew(t *testing.T) {
	cache := newMapItemGetterSetter()
	db := newMapItemGetterSetter()
	_, err := New(cache, db)
	if err != nil {
		t.Fatalf("ItemService.New: %v\n", err)
	}
}

func TestSet(t *testing.T) {
	cache := newMapItemGetterSetter()
	db := newMapItemGetterSetter()
	itemService, err := New(cache, db)
	if err != nil {
		t.Fatalf("ItemService.New: %v\n", err)
	}
	want := Item{Key: "key", Value: "value"}
	// set item in service
	// should propagate to db and cache
	err = itemService.Set(want)
	if err != nil {
		t.Fatalf("ItemService.Set: %v\n", err)
	}
	// should be in cache
	got, err := cache.Get(want.Key)
	if err != nil {
		t.Fatalf("cache.Get: %v\n", err)
	}
	if got != want {
		t.Fatalf("wanted %#v but got %#v\n", want, got)
	}
	// should be in db
	got, err = db.Get(want.Key)
	if err != nil {
		t.Fatalf("db.Get: %v\n", err)
	}
	if got != want {
		t.Fatalf("wanted %#v but got %#v\n", want, got)
	}
}

func TestGet(t *testing.T) {
	cache := newMapItemGetterSetter()
	db := newMapItemGetterSetter()
	itemService, err := New(cache, db)
	if err != nil {
		t.Fatalf("ItemService.New: %v\n", err)
	}
	want := Item{Key: "key", Value: "value"}
	// set item in service
	// should propagate to db and cache
	err = itemService.Set(want)
	if err != nil {
		t.Fatalf("ItemService.Set: %v\n", err)
	}

	got, err := itemService.Get(want.Key)
	if err != nil {
		t.Fatalf("ItemService.Get: %v\n", err)
	}
	if got != want {
		t.Fatalf("wanted %#v but got %#v\n", want, got)
	}
}

func TestGetCache(t *testing.T) {
	cache := newMapItemGetterSetter()
	db := newMapItemGetterSetter()
	itemService, err := New(cache, db)
	if err != nil {
		t.Fatalf("ItemService.New: %v\n", err)
	}
	want := Item{Key: "key", Value: "value"}
	// set item in db
	err = db.Set(want)
	if err != nil {
		t.Fatalf("db.Set: %v\n", err)
	}

	// should get from db
	got, err := itemService.Get(want.Key)
	if err != nil {
		t.Fatalf("ItemService.Get: %v\n", err)
	}
	if got != want {
		t.Fatalf("wanted %#v but got %#v\n", want, got)
	}

	// should populate cache
	got, err = cache.Get(want.Key)
	if err != nil {
		t.Fatalf("cache.Get: %v\n", err)
	}
	if got != want {
		t.Fatalf("wanted %#v but got %#v\n", want, got)
	}
}

func TestDelete(t *testing.T) {
	cache := newMapItemGetterSetter()
	db := newMapItemGetterSetter()
	itemService, err := New(cache, db)
	if err != nil {
		t.Fatalf("ItemService.New: %v\n", err)
	}
	want := Item{Key: "key", Value: "value"}
	// set item in service
	// should propagate to db and cache
	err = itemService.Set(want)
	if err != nil {
		t.Fatalf("ItemService.Set: %v\n", err)
	}

	got, err := itemService.Get(want.Key)
	if err != nil {
		t.Fatalf("ItemService.Get: %v\n", err)
	}
	if got != want {
		t.Fatalf("wanted %#v but got %#v\n", want, got)
	}

	err = itemService.Delete(want.Key)
	if err != nil {
		t.Fatalf("ItemService.Set: %v\n", err)
	}

	got, err = itemService.Get(want.Key)
	fmt.Println(err.Error())
	fmt.Println(reflect.TypeOf(err))
	if err == nil {
		t.Fatalf("ItemService.Get: %v\n", "item should be missing")
	}
}
