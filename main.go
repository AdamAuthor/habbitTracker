package main

import (
	"awesomeProject/server/db/newdb"
	"awesomeProject/server/router"
	"context"
	"fmt"
	"log"
	"os"
)

func main() {
	url := os.Getenv("DATABASE_URL")

	database := newdb.NewDB()
	if err := database.Connect(url); err != nil {
		panic(err)
	}
	defer func() {
		err := database.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//Creating and run new server
	srv := router.NewServer(context.Background(), ":"+port, database)
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}

	srv.WaitForGT()

}
