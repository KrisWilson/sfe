package listener

import (
	_ "database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sfe/settings"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type UserJSON struct {
	Name string `json:"user"`
	Pass string `json:"pass"`
}

func exploreHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		var u = CheckToken(r.Header.Get("Token"))
		if u.ID != -1 {
			//		_, err := fmt.Fprintf(w, "Authorized - token accepted, Welcome "+u.Name+"\n")
			//		if err != nil {
			//			return
			//		}

			err := r.ParseForm()
			if err != nil {
				return
			}
			now := time.Now()
			// TODO: Uodpornić parametr "path" na exploracje całej przestrzenii dyskowej tj. "../../../"
			// TODO: Zgłębić temat czy wykonywanie os.Mkdir etc. zezwala na overbuffer/injection do shella
			folderPath := settings.Load().SharedDir + r.FormValue("path")
			folderPath = strings.Replace(folderPath, "..", ".", -1)
			file := r.FormValue("file")
			file = strings.Replace(file, "..", ".", -1)
			// TODO: Rozdzielić logikę explorer od download file
			// TODO: Dodać download dir
			// TODO: Dodać wielewątków TCP w celu szybszego pobierania danych oraz weryfikacje pobierania danych
			// TODO: poprawić logikę pobierania plików

			err = os.Mkdir(settings.Load().SharedDir, 0777)
			if err != nil {
			} else {
				err := os.Mkdir(settings.Load().SharedDir+"/Pics", 0777)
				if err != nil {
					return
				}
				err = os.Mkdir(settings.Load().SharedDir+"/Files", 0777)
				if err != nil {
					return
				}
				err = os.WriteFile(settings.Load().SharedDir+"/some.file", []byte("some.file's content <<>>!...!<<>>"), 0777)
				if err != nil {
					return
				}
			}

			if file != "" {
				if len(folderPath) == 0 {
					folderPath = "/"
				}
				fmt.Println(now.Format(time.DateTime) + " [Explorer] " + u.Name + " accessed file: " + folderPath + file + "\r")
				fileDownload, _ := os.ReadFile(folderPath + file)
				_, err := w.Write(fileDownload)
				if err != nil {
					return
				}
			} else {
				files, err := os.ReadDir(folderPath)
				if err != nil {
					fmt.Println("Error reading directory:\r", err, "\r")
					_, err := fmt.Fprintln(w, "Error reading directory "+folderPath+"\r")
					if err != nil {
						return
					}
					return
				}
				fmt.Println(now.Format(time.DateTime) + " [Explorer] " + u.Name + " accessed folder: " + folderPath + "\r")
				_, err = fmt.Fprintf(w, "Folder path: %s\n\r", folderPath)
				if err != nil {
					return
				}
				for _, file := range files {
					if file.IsDir() {
						_, err2 := fmt.Fprintf(w, "Folder\t"+file.Name()+"\n\r")
						if err2 != nil {
							return
						}
					} else {
						_, err := fmt.Fprintf(w, "File\t"+file.Name()+"\n\r")
						if err != nil {
							return
						}
					}
				}
			}
		} else {
			_, err := fmt.Fprintf(w, "Authorized - token not accepted"+"\r")
			if err != nil {
				return
			}
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
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var user UserJSON
	err = json.Unmarshal(bodyBytes, &user)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	if CheckPassword(user.Pass, user.Name) {
		//		_, err := fmt.Fprintln(w, "<<Token::Authorized>>")
		_, err = fmt.Fprint(w, newToken(user.Name))
		now := time.Now()
		fmt.Println(now.Format(time.DateTime) + " [authorizeHandler] " + r.URL.Path + " " + r.Method + " " + r.RemoteAddr + " Authorized: " + user.Name + "\r")
		if err != nil {
			return
		}
	} else {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println("[authorizeHandler] " + r.URL.Path + " " + r.Method + " " + r.RemoteAddr + " Unauthorized: " + user.Name + "\r")
		return
	}

}

func Host(port int) {
	config := settings.Load()
	InitDB(config.ServerDB)

	http.HandleFunc("/authorize", authorizeHandler)
	http.HandleFunc("/explore", exploreHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {})

	if port != -1 {
		log.Printf("Listening on port %d\n\r", port)
		err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
		if err != nil {
			panic(err)
		}
	} else {
		log.Printf("Listening on port %d\n\r", config.ServerPort)
		err := http.ListenAndServe(":"+strconv.Itoa(config.ServerPort), nil)
		if err != nil {
			panic(err)
		}
	}

}
