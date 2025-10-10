package main

import (
	"fmt"

	"github.com/sharifahmad2061/trip-grpc-go/internal/config"
)

func main() {
	conf := config.Load()
	fmt.Println(conf)
}
