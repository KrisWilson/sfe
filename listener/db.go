package listener

import (
	"bufio"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"sfe/settings"
	"strings"
	"time"
)

type User struct {
	ID      int       `db:"id"`
	Name    string    `db:"name"`
	Pass    string    `db:"pass"`
	Dir     string    `db:"dir"`
	Token   string    `db:"token"`
	Timeout time.Time `db:"timeout"`
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

		fmt.Println("[DB] Utworzono baze danych: " + dbname + ".db")

		dsn := "file:" + dbname + ".db?cache=shared&mode=rw"
		db, err := sql.Open("sqlite3", dsn)
		if err != nil {
			fmt.Println("[DB] Unable to connect to database")
			panic(err)
		}

		defer func(db *sql.DB) {
			err := db.Close()
			if err != nil {
			}
		}(db)

		_, err = db.Exec(
			`CREATE TABLE IF NOT EXISTS users (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				name VARCHAR(20) NOT NULL,
				pass VARCHAR(20) NOT NULL,
				dir TEXT NOT NULL,
				token CHAR(64) UNIQUE NOT NULL,
    			timeout DATETIME NOT NULL 
                )`)
		if err != nil {
			fmt.Println("[DB] Unable to create table")
			panic(err)
		}
		fmt.Println("[DB] Dodano tablice users do bazy danych")
		err = db.Close()
		if err != nil {
			fmt.Println("[DB] Unable to close database")
			return
		}
		AddUser("user", "password", "userdir")
	}
}

func newToken(username string) string {
	dsn := "file:" + settings.Load().ServerDB + ".db?cache=shared&mode=rw"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		fmt.Println("[DB] Unable to connect to database")
	}
	var token string
	for i := 0; i < 64; i++ {
		token = token + string(rune('A'+rand.Intn(52)))
	}

	err = db.QueryRow("UPDATE users SET token = ? WHERE name = ?", token, username).Scan(&username)
	// dodajemy time.Hour*1 - raz bo jesteśmy w strefie +01:00 oraz drugi raz ponieważ token będzie ważny tyle
	err = db.QueryRow("UPDATE users SET timeout = ? WHERE name = ?", time.Now().Add(time.Hour*1+time.Hour*1).Format(time.DateTime), username).Scan(&username)
	return token
}

func getUser(username string) (User, error) {

	dsn := "file:" + settings.Load().ServerDB + ".db?cache=shared&mode=rw"
	db, err := sql.Open("sqlite3", dsn)

	var u User
	err = db.QueryRow("SELECT * FROM users WHERE name = ?", username).Scan(&u.ID, &u.Name, &u.Pass, &u.Dir, &u.Token, &u.Timeout)
	if err != nil {
		return u, err
	}
	return u, nil
}

func RemoveUser(username string) {
	dsn := "file:" + settings.Load().ServerDB + ".db?cache=shared&mode=rw"
	db, _ := sql.Open("sqlite3", dsn)
	_, err := db.Exec("DELETE FROM users WHERE name = ?", username)
	if err != nil {
		fmt.Println("[DB] Unable to remove user")
	}
}

func AddUser(username string, password string, userdir string) {
	dsn := "file:" + settings.Load().ServerDB + ".db?cache=shared&mode=rw"
	db, err := sql.Open("sqlite3", dsn)
	if err != nil {
		fmt.Println("[DB] Unable to connect to database")
	}

	if len(userdir) == 0 {
		userdir = username + "dir"
	}

	now := time.Now()
	_, err = db.Exec("INSERT INTO users (name, pass, dir, token, timeout) VALUES (?,?,?,?,?)",
		username, password, userdir, "", now.Format(time.DateTime))
	if err != nil {
		fmt.Println("[DB] Unable to insert user")
		panic(err)
	}
	// Nowy token - ustawienie tokenu DB na unique, nie mogą być 2 puste pola bo to nie są unikalne
	newToken(username)

	fmt.Println("Dodano uzytkownika " + username)
}

// CheckPassword - Sprawdzenie poprawności hasła użytkownika który się zgłasza w celu potwierdzenia autoryzacji
func CheckPassword(password string, user string) bool {
	u, err := getUser(user)
	if err != nil {
		fmt.Println("Wystąpił problem: " + user + " => " + err.Error())
		return false
	}
	return u.Pass == password
}

func CheckToken(token string) User {
	dsn := "file:" + settings.Load().ServerDB + ".db?cache=shared&mode=rw"
	db, err := sql.Open("sqlite3", dsn)
	var u User
	u.ID = -1
	//fmt.Println(token)
	err = db.QueryRow("SELECT * FROM users WHERE token = ?", token).Scan(&u.ID, &u.Name, &u.Pass, &u.Dir, &u.Token, &u.Timeout)
	if err != nil {
		fmt.Println("Token not found")
		return u
	}
	//fmt.Println(u.Token)
	now := time.Now()
	//fmt.Println("Now: " + now.Format(time.RFC822Z) + " :: Token: " + u.Timeout.Format(time.RFC822Z))
	//if u.Timeout.After(now) { // czy timeout jest po aktualnym czasie tzn później
	if now.After(u.Timeout) {
		fmt.Println("[DB] Token expired")
		u.ID = -1
		return u
	}
	return u
}

func readKey() rune {
	reader := bufio.NewReader(os.Stdin)
	char, _, _ := reader.ReadRune()
	return char
}

func ConfigDB() {
	InitDB(settings.Load().ServerDB)
	fmt.Println("\033[31m<<< \u001B[0mSFE - DB Config SFE \u001B[31m>>>\u001B[0m\r")
	fmt.Println("[1] Show users from sqlite database\r")
	fmt.Println("[2] Add user\r")
	fmt.Println("[3] Remove user\r")
	fmt.Println("[X] Exit\r")
	fmt.Print("Your choice: ")
	input := readKey()

	switch string(input) {
	case "1": // pokazanie wszystkich uzytkowników
		dsn := "file:" + settings.Load().ServerDB + ".db?cache=shared&mode=rw"
		db, _ := sql.Open("sqlite3", dsn)

		rows, err := db.Query("SELECT * FROM users")
		if err != nil {
			panic(err.Error())
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {

			}
		}(rows)

		fmt.Println("ID \t| Username \t| Password \t| Dir \t| Token \t| Timeout")
		for rows.Next() {
			var col1, col2, col3, col4, col5 string
			var col6 time.Time
			if err := rows.Scan(&col1, &col2, &col3, &col4, &col5, &col6); err != nil {
				panic(err.Error())
			}
			fmt.Printf("%s \t| %s \t| %s \t| %s \t| %s \t| %s\n", col1, col2, col3, col4, col5, col6.Format(time.DateTime))
		}

		if err := rows.Err(); err != nil {
			panic(err.Error())
		}

	case "2": // dodanie użytkownika
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Podaj nazwe uzytkownika: ")
		username, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Błąd podczas czytania danych:", err)
			return
		}

		fmt.Print("Podaj hasło uzytkownika: ")
		password, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Błąd podczas czytania danych:", err)
			return
		}

		fmt.Print("Podaj folder uzytkownika [default=username]: ")
		userdir, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Błąd podczas czytania danych:", err)
			return
		}
		username = strings.ReplaceAll(username, "\n", "")
		password = strings.ReplaceAll(password, "\n", "")
		userdir = strings.ReplaceAll(userdir, "\n", "")
		AddUser(username, password, userdir)

	case "3":
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Podaj nazwe uzytkownika: ")
		username, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Błąd podczas czytania danych:", err)
			return
		}
		username = strings.ReplaceAll(username, "\n", "")
		RemoveUser(username)
	case "X":
		return

	default:
		fmt.Println("Invalid choice")
		fmt.Print("<< Press enter to continue\n")
		reader := bufio.NewReader(os.Stdin)
		_, _, _ = reader.ReadRune()
		ConfigDB()
	}

}
