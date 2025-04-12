package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"go_final_project/pkg/api"
	"go_final_project/tests"
)

func Run(ts api.TaskService) {
	webDir := "./web"
	port, err := strconv.Atoi(os.Getenv("TODO_PORT"))
	if err != nil {
		port = tests.Port
	}
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	api.Init(ts)

	fmt.Printf("Server running on port %d...", port)
	err = http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

}
