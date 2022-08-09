package cache

import (
	"fmt"
	"testing"
	"time"
)

func TestItemAdds(t *testing.T) {
	instance := NewTTLMap(10, time.Hour)

	key := "key1"
	value := "value1"

	instance.Add(key, value)

	data, err := instance.Get(key)
	if err != nil {
		t.Fatal(err)
	}

	if data != value {
		t.Fatalf("Expected %s, got %s", value, data)
	}
}

func TestItemDeletes(t *testing.T) {
	instance := NewTTLMap(10, time.Hour)

	key := "key2"
	value := "value2"

	instance.Add(key, value)
	instance.Delete(key)

	_, err := instance.Get(key)
	if err == nil {
		t.Fatalf("Expected item to have been deleted")
	}
}

func TestItemDoesExpire(t *testing.T) {
	instance := NewTTLMap(10, time.Second)

	key := "key3"
	value := "value3"

	instance.Add(key, value)

	time.Sleep(time.Second * 3)

	_, err := instance.Get(key)
	if err == nil {
		t.Fatalf("Expected item to have expired")
	}
}

func TestMaxItems(t *testing.T) {
	instance := NewTTLMap(3, time.Hour)
	for i := 1; i <= 10; i++ {
		instance.Add(fmt.Sprintf("key%d", i), "value")
	}

	if len(instance.m) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(instance.m))
	}
}

func TestMaxItemsEvictsOldest(t *testing.T) {
	instance := NewTTLMap(3, time.Hour)
	for i := 1; i <= 10; i++ {
		instance.Add(fmt.Sprintf("key%d", i), "value")
	}

	if _, err := instance.Get("key1"); err == nil {
		t.Fatalf("Expected item to have been evicted")
	}

	if _, err := instance.Get("key7"); err == nil {
		t.Fatalf("Expected item to have been evicted")
	}

	if _, err := instance.Get("key8"); err != nil {
		t.Fatalf("Expected item to not have been evicted")
	}

	if _, err := instance.Get("key9"); err != nil {
		t.Fatalf("Expected item to not have been evicted")
	}

	if _, err := instance.Get("key10"); err != nil {
		t.Fatalf("Expected item to not have been evicted")
	}
}
