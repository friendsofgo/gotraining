package main

import (
	"fmt"

	"github.com/google/uuid"
)

type RedisFake struct {
	ID string
}

func NewRedisFake() *RedisFake {
	return &RedisFake{uuid.New().String()}
}

func (r *RedisFake) Close() error{
	fmt.Printf("closing resource: %s...\n", r.ID)
	return nil
}