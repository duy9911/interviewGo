package main

import (
	"interview1710/api/routers"
)

func main() {
	// seed.CreateTable()
	// seed.SeedData()
	// elasticDB.AddToEs()
	routers.HandleRequests()
	// ExampleClient()
}
