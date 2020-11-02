package main

import (
	"context"
	"io"
	"log"
	"math/rand"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	createConnection := func() (io.Closer, error) {
		return NewRedisFake(), nil
	}

	pool, err := NewPool(2, createConnection)

	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	g, _ := errgroup.WithContext(context.Background())

	for i := 0; i < 25; i++ {
		i := i
		g.Go(func() error {
			return simulateAction(i, pool)
		})
	}
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}
}

func simulateAction(query int, pool *Pool) error {
	conn, err := pool.Get()
	if err != nil {
		return err
	}
	defer pool.Release(conn)

	// Simulate time to the action
	time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)
	log.Printf("Query ID[%d] Connection ID[%s]\n", query, conn.(*RedisFake).ID)
	return nil
}
