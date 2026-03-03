package listener

import (
	"log"
	"net/http"
	"strconv"
)

func Host(port int) {
	log.Printf("Listening on port %d\n", port)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	err := http.ListenAndServe(":"+strconv.Itoa(port), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	if err != nil {
		panic(err)
	}
}
