package main

import (
	"context"
	"log"

	"distributed-kvs/internal/app/node"
)

func main() {
	err := node.Run(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
}
