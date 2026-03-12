package listener

import (
	"bufio"
	"database/sql"
	"fmt"
	"math/rand"
	"os"
	"sfe/settings"
)

type user struct {
	ID    int    `db:"id"`
	Name  string `db:"name"`
	Pass  string `db:"pass"`
	Dir   string `db:"dir"`
	Token string `db:"token"`
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
				dir TEXT,
				token CHAR(64) UNIQUE
                )`)
		if err != nil {
			fmt.Println("[DB] Unable to create table")
			panic(err)
		}
		fmt.Println("[DB] Dodano tablice users do bazy danych")

		_, err = db.Exec("INSERT INTO users (name, pass, dir, token) VALUES (?,?,?,?)", "user", "password", "userdir", "")
		if err != nil {
			fmt.Println("[DB] Unable to insert user")
			panic(err)
		}

		fmt.Println("[DB] Dodano użytkownika defaultowego user:password")
		err = db.Close()
		if err != nil {
			fmt.Println("[DB] Unable to close database")
			return
		}
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
		token = token + string('A'+rand.Intn(52))
	}

	err = db.QueryRow("UPDATE users SET token = ? WHERE name = ?", token, username).Scan(&username)
	return token
}

func getUser(username string) (user, error) {

	dsn := "file:" + settings.Load().ServerDB + ".db?cache=shared&mode=rw"
	db, err := sql.Open("sqlite3", dsn)

	var u user
	err = db.QueryRow("SELECT * FROM users WHERE name = ?", username).Scan(&u.ID, &u.Name, &u.Pass, &u.Dir, &u.Token)
	if err != nil {
		return u, err
	}
	return u, nil
}

func RemoveUser(username string) {
	dsn := "file:" + settings.Load().ServerDB + ".db?cache=shared&mode=rw"
	db, _ := sql.Open("sqlite3", dsn)
	db.QueryRow("DELETE FROM users WHERE name = ?", username)
}

func AddUser(username string, password string, userdir string) {
	dsn := "file:" + settings.Load().ServerDB + ".db?cache=shared&mode=rw"
	db, _ := sql.Open("sqlite3", dsn)
	if len(userdir) == 0 {
		userdir = username + "dir"
	}
	db.QueryRow("INSERT INTO users (name, pass, dir) VALUES (?, ?, ?)", username, password, userdir)
	fmt.Println("Dodano uzytkownika " + username)
}

// Sprawdzenie poprawności hasła użytkownika który się zgłasza w celu potwierdzenia autoryzacji
func CheckPassword(password string, user string) bool {
	u, err := getUser(user)
	if err != nil {
		fmt.Println("Wystąpił problem: " + user + " => " + err.Error())
		return false
	}
	return u.Pass == password
}

func CheckToken(token string) user {
	dsn := "file:" + settings.Load().ServerDB + ".db?cache=shared&mode=rw"
	db, err := sql.Open("sqlite3", dsn)
	var u user
	u.ID = -1
	err = db.QueryRow("SELECT * FROM users WHERE token = ?", token).Scan(&u.ID, &u.Name, &u.Pass, &u.Dir, &u.Token)
	if err != nil {
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
		defer rows.Close()

		fmt.Println("ID \t| Username \t| Password \t| Dir")
		for rows.Next() {
			var col1, col2, col3, col4 string

			if err := rows.Scan(&col1, &col2, &col3, &col4); err != nil {
				panic(err.Error())
			}
			fmt.Printf("%s \t| %s \t| %s \t| %s\n", col1, col2, col3, col4)
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
		AddUser(username, password, userdir)

	case "3":

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
