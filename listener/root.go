package listener

import (
	_ "database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sfe/settings"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type User struct {
	Name string `json:"user"`
	Pass string `json:"pass"`
}

func exploreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var u = CheckToken(r.Header.Get("Token"))
		if u.ID != -1 {
			_, err := fmt.Fprintf(w, "Authorized - token accepted, Welcome "+u.Name+"\n")
			if err != nil {
				return
			}

			err = r.ParseForm()
			if err != nil {
				return
			}
			now := time.Now()
			// TODO: Uodpornić parametr "path" na exploracje całej przestrzenii dyskowej tj. "../../../"
			folderPath := settings.Load().Shared + r.FormValue("path")
			file := r.FormValue("file")
			// TODO: Rozdzielić logikę explorer od download file
			// TODO: Dodać download dir
			// TODO: Dodać wielewątków TCP w celu szybszego pobierania danych oraz weryfikacje pobierania danych
			if file != "" {
				if len(folderPath) == 0 {
					folderPath = "/"
				}
				fmt.Println(now.Format(time.DateTime) + " [Explorer] " + u.Name + " accessed file: " + folderPath + file)
				fileDownload, _ := os.ReadFile(folderPath + file)
				_, err := w.Write(fileDownload)
				if err != nil {
					return
				}
			} else {
				files, err := os.ReadDir(folderPath)
				if err != nil {
					fmt.Println("Error reading directory:", err)
					fmt.Fprintln(w, "Error reading directory")
					return
				}
				fmt.Println(now.Format(time.DateTime) + " [Explorer] " + u.Name + " accessed folder: " + folderPath)
				fmt.Fprintf(w, "Folder path: %s\n", folderPath)
				for _, file := range files {
					if file.IsDir() {
						fmt.Fprintf(w, "Folder\t"+file.Name()+"\n")
					} else {
						fmt.Fprintf(w, "File\t"+file.Name()+"\n")
					}
				}
			}
		} else {
			fmt.Fprintf(w, "Authorized - token not accepted")
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func authorizeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
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

	if CheckPassword(user.Pass, user.Name) {
		//		_, err := fmt.Fprintln(w, "<<Token::Authorized>>")
		_, err = fmt.Fprint(w, newToken(user.Name))
		now := time.Now()
		fmt.Println(now.Format(time.DateTime) + " [authorizeHandler] " + r.URL.Path + " " + r.Method + " " + r.RemoteAddr + " Authorized: " + user.Name)
		if err != nil {
			return
		}
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println("[authorizeHandler] " + r.URL.Path + " " + r.Method + " " + r.RemoteAddr + " Unauthorized: " + user.Name)
		return
	}

}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintln(w, "List of users...")
	if err != nil {
		return
	}
}

func Host(port int) {
	config := settings.Load()
	InitDB(config.ServerDB)

	http.HandleFunc("/authorize", authorizeHandler)
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/explore", exploreHandler)

	log.Printf("Listening on port %d\n", config.ServerPort)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	err := http.ListenAndServe(":"+strconv.Itoa(config.ServerPort), nil)
	if err != nil {
		panic(err)
	}
}
