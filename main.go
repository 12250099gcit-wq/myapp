// myapp/main.go
package main

import (
	"log"

	"myapp/routes"
	"myapp/utils/postgres"
)

func main() {
	if err := postgres.Init(); err != nil {
		log.Fatal(err)
	}

	routes.StartServer()
}

