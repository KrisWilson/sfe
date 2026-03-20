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

type FileJSON struct {
	Name         string `json:"Name"`
	Type         string `json:"Type"`
	Size         int64  `json:"Size"`
	DateModified string `json:"DateModified"`
}

var config settings.Config

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "PUT" {
		var u = CheckToken(r.Header.Get("Token"))
		if u.ID != -1 {
			err := r.ParseForm()
			if err != nil {
				return
			}

			filename := r.FormValue("filename")
			uploadPath := r.FormValue("uploadpath")

			// take off "../../../" from filename and uploadpath
			i := 0
			for i < 1 {
				if strings.Contains(filename, "..") {
					filename = strings.Replace(filename, "..", ".", -1)
				} else {
					i = 2
				}
			}

			i = 0
			for i < 1 {
				if strings.Contains(uploadPath, "..") {
					uploadPath = strings.Replace(uploadPath, "..", ".", -1)
				} else {
					i = 2
				}
			}

			// utworz folder uploadu danego uzytkownika
			err = os.MkdirAll(config.SharedDir+"/"+u.Dir+"/"+uploadPath, os.ModePerm)
			if err != nil {
				_, err := fmt.Fprintln(w, err)
				if err != nil {
					return
				}
				return
			}
			if filename == "." {
				_, err = fmt.Fprint(w, "[Server Feedback] \u001B[36mFolder "+uploadPath+" craeted successfully\r")
				fmt.Println("[Uploader] " + u.Name + " has saved " + u.Dir + "/" + uploadPath + "\r")
				return
			}
			respBody, err := io.ReadAll(r.Body)
			if err != nil {
				_, err = fmt.Fprintln(w, err)
			}
			err = os.WriteFile(config.SharedDir+"/"+u.Dir+"/"+uploadPath+"/"+filename, respBody, os.ModePerm)
			if err != nil {
				_, err = fmt.Fprintln(w, err)
				return
			} else {
				_, err = fmt.Fprint(w, "[Server Feedback] \u001B[36mFile "+filename+" uploaded successfully\r")
				fmt.Println("[Uploader] " + u.Name + " has saved " + u.Dir + uploadPath + "/" + filename + "\r")
			}
		} else {
			http.Error(w, "Forbidden; Wrong Token", http.StatusForbidden)
			return
		}
	} else {
		http.Error(w, "Method not allowed; Only PUT is allowed", http.StatusMethodNotAllowed)
		return
	}
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
				fmt.Println(err)
				return
			}
			now := time.Now() //?
			// TODO: Zgłębić temat czy wykonywanie os.Mkdir etc. zezwala na overbuffer/injection do shella
			folderPath := config.SharedDir + r.FormValue("path")

			i := 0
			for i < 1 {
				if strings.Contains(folderPath, "..") {
					folderPath = strings.Replace(folderPath, "..", ".", -1)
				} else {
					i = 2
				}
			}

			file := r.FormValue("file")

			i = 0
			for i < 1 {
				if strings.Contains(file, "..") {
					file = strings.Replace(file, "..", ".", -1)
				} else {
					i = 2
				}
			}

			//fmt.Println("Plik: " + file + "\r")
			//fmt.Println("Folder: " + folderPath + "\r")
			// TODO: Rozdzielić logikę explorer od download file
			// TODO: Dodać wielewątków TCP w celu szybszego pobierania danych oraz weryfikacje pobierania danych
			// TODO: poprawić logikę pobierania plików

			// przygotowanie przykładowych folderów do użycia, aby zapewnić istnienie folderu /share
			err = os.Mkdir(config.SharedDir, 0777)
			if err != nil {
			} else {
				err := os.Mkdir(config.SharedDir+"/Pics", 0777)
				if err != nil {
					return
				}
				err = os.Mkdir(config.SharedDir+"/Files", 0777)
				if err != nil {
					return
				}
				err = os.WriteFile(config.SharedDir+"/some.file", []byte("some.file's content <<>>!...!<<>>"), 0777)
				if err != nil {
					return
				}
			}

			// Zdefiniowany plik - wystawiamy plik do pobrania
			if file != "" {
				if len(folderPath) == 0 {
					folderPath = "/"
				}

				fileDownload, errfile := os.ReadFile(folderPath + file)
				if errfile != nil {
					fmt.Println(now.Format(time.DateTime) + " [Explorer] " + u.Name + " CAN'T access file: " + folderPath + file + "\r")
					w.WriteHeader(http.StatusNotFound)
					fmt.Fprint(w, errfile.Error())
					return
				} else {
					w.WriteHeader(http.StatusOK)
					fmt.Println(now.Format(time.DateTime) + " [Explorer] " + u.Name + " accessed file: " + folderPath + file + "\r")
				}
				_, err := w.Write(fileDownload)
				if err != nil {
					return
				}
			} else { // Brak zdefiniowanego pliku - tutaj jest logika listingu plików w folderze
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
				//_, err = fmt.Fprintf(w, "Folder path: %s\n\r", folderPath)
				//if err != nil {
				//	return
				//}

				var filesJson []FileJSON

				for _, file := range files {

					fileInfo, _ := os.Stat(folderPath + file.Name())
					if file.IsDir() {
						filesJson = append(filesJson, FileJSON{
							Name:         fileInfo.Name(),
							Type:         "Folder",
							Size:         0,
							DateModified: fileInfo.ModTime().Format("2006-01-02 15:04:05"),
						})
					} else {
						filesJson = append(filesJson, FileJSON{
							Name:         fileInfo.Name(),
							Type:         "File",
							Size:         fileInfo.Size(),
							DateModified: fileInfo.ModTime().Format("2006-01-02 15:04:05"),
						})
					}
				}

				encoder := json.NewEncoder(w)
				err2 := encoder.Encode(&filesJson)
				if err2 != nil {
					http.Error(w, "Internal Server Error", 500)
					return
				}
				w.Header().Set("Content-Type", "application/json")

			}
		} else {
			fmt.Println()
			http.Error(w, "Authorized - token not accepted", http.StatusBadRequest)
			return
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
		fmt.Println("Authorizing: " + user.Pass + " : " + user.Name + "\r")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		fmt.Println("[authorizeHandler] " + r.URL.Path + " " + r.Method + " " + r.RemoteAddr + " Unauthorized: " + user.Name + "\r")
		return
	}

}

func Host(port int) {
	config = settings.Load()
	InitDB(config.ServerDB)

	http.HandleFunc("/authorize", authorizeHandler)
	http.HandleFunc("/explore", exploreHandler)
	http.HandleFunc("/upload", uploadHandler)
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
