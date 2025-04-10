package main

import (
    "fmt"
    "go_final_project/pkg/db"
    "go_final_project/pkg/server"
)

func main() {
    err := db.Init()
    if err != nil {
        fmt.Println(err)
    }
    server.Run()
}
