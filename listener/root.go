package listener

import (
	"encoding/json"
	_ "encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sfe/settings"
	"strconv"
)

type User struct {
	Pass string `json:"pass"`
}

func authorizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var user User
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	fmt.Println("[authorizeHandler] " + r.URL.Path + " " + r.Method + " " + r.RemoteAddr + " " + r.FormValue("pass") + user.Pass)
	if settings.PassVerify(user.Pass) {
		_, err := fmt.Fprintln(w, "Authorized")
		if err != nil {
			return
		}
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
