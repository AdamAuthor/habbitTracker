package main

import (
	"awesomeProject/server/db/newdb"
	"awesomeProject/server/router"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	url := os.Getenv("URL")

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

	//Creating and run new server
	srv := router.NewServer(context.Background(), "https://habitadam.herokuapp.com/", database)
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}

	srv.WaitForGT()

}
