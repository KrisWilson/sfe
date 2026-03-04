package listener

import (
	"database/sql"
	"fmt"
	"os"
	"sfe/settings"
)

type user struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
	Pass string `db:"pass"`
}

func InitDB(dbname string) {
	_, err := os.Stat(dbname + ".db")
	if os.IsNotExist(err) {
		file, err := os.Create(dbname + ".db")
		if err != nil {
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {

			}
		}(file)

		fmt.Println("Utworzono baze danych: " + dbname + ".db")

		dsn := "file:" + dbname + ".db?cache=shared&mode=rw"
		db, err := sql.Open("sqlite3", dsn)
		if err != nil {
			fmt.Println("Unable to connect to database")
			panic(err)
		}

		defer func(db *sql.DB) {
			fmt.Println("!")
			err := db.Close()
			if err != nil {
			}
		}(db)

		_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
                id INTEGER PRIMARY KEY,
                name TEXT NOT NULL,
                pass TEXT NOT NULL)`)
		if err != nil {
			fmt.Println("Unable to create table")
			panic(err)
		}
		fmt.Println("Dodano tablice users do bazy danych")

		_, err = db.Exec("INSERT INTO users (id, name, pass) VALUES (?, ?, ?)", 1, "user", "password")
		if err != nil {
			fmt.Println("Unable to insert user")
			panic(err)
		}

		fmt.Println("Dodano użytkownika defaultowego user:password")
		err = db.Close()
		if err != nil {
			fmt.Println("Unable to close database")
			return
		}
	}
}

func getUser(username string) (user, error) {

	dsn := "file:" + settings.Load().ServerDB + ".db?cache=shared&mode=rw"
	db, err := sql.Open("sqlite3", dsn)

	var u user
	err = db.QueryRow("SELECT * FROM users WHERE name = ?", username).Scan(&u.ID, &u.Name, &u.Pass)
	if err != nil {
		return u, err
	}
	return u, nil
}

// Sprawdzenie poprawności hasła
func CheckPassword(password string, user string) bool {
	u, err := getUser(user)
	if err != nil {
		panic(err)
	}
	return u.Pass == password
}
