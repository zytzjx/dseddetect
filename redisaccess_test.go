package main

import (
	"fmt"
	"testing"
)

func TestCreateRedisPool(t *testing.T) {
	CreateRedisPool(5)
	ping(1)
}
func TestPublishData(t *testing.T) {
	CreateRedisPool(48)
	a, err := Publish(6, "status", "connected")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(a)
}
