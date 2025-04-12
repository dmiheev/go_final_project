package main

import (
	"fmt"
	"go_final_project/pkg/api"
	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
)

func main() {

	dbConn, err := db.Init()
	if err != nil {
		fmt.Println(err)
	}
	defer dbConn.Close()

	storage := db.NewStorage(dbConn)
	service := api.NewTaskService(storage)

	server.Run(service)
}
