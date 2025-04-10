package server

import (
    "go_final_project/pkg/api"
    "go_final_project/tests"
    "log"
    "net/http"
    "os"
    "strconv"
)

func Run() {
    webDir := "./web"
    port, err := strconv.Atoi(os.Getenv("TODO_PORT"))
    if err != nil {
        port = tests.Port
    }
    http.Handle("/", http.FileServer(http.Dir(webDir)))

    api.Init()

    err = http.ListenAndServe(":"+strconv.Itoa(port), nil)
    if err != nil {
        log.Fatal("ListenAndServe: ", err)
    }
}
