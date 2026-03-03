package listener

import (
	"fmt"
	"log"
	"net/http"
	"sfe/settings"
	"strconv"
)

func authorizeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("[authorizeHandler] " + r.URL.Path + " " + r.Method + " " + r.RemoteAddr)
	fmt.Println("! any")
	_, err := fmt.Fprintln(w, "Hello, world!")
	if err != nil {
		return
	}
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "List of users...")
}

func Host(port int) {

	config := settings.Load()
	http.HandleFunc("/authorize", authorizeHandler)
	http.HandleFunc("/users", usersHandler)
	log.Printf("Listening on port %d\n", config.ServerPort)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	err := http.ListenAndServe(":"+strconv.Itoa(config.ServerPort), nil)
	if err != nil {
		panic(err)
	}
}
