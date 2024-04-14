package main

import (
	"context"
	"log"

	"distributed-kvs/internal/app/cluster"
)

func main() {
	err := cluster.Run(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
}
