package client

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"sfe/listener"
	"sfe/settings"

	"github.com/spf13/viper"
)

func readKey() rune {
	reader := bufio.NewReader(os.Stdin)
	char, _, _ := reader.ReadRune()
	return char
}

func Run() {

	//// switch stdin into 'raw' mode
	//oldState, err := term.MakeRaw(in t(os.Stdin.Fd()))
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//defer term.Restore(int(os.Stdin.Fd()), oldState)

	//b := make([]byte, 1)
	//_, err = os.Stdin.Read(b)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Printf("the char %q was hit", string(b[0]))

	fmt.Println("\033[31m<<< \u001B[0mSFE - Small File Exchanger \u001B[31m>>>\u001B[0m\r")
	fmt.Println("[1] Connect to Server\r")
	fmt.Println("[2] Host a server\r")
	fmt.Println("[3] Show config\r")
	fmt.Println("[4] Config DB\r")
	fmt.Println("[X] Exit\r")
	fmt.Print("Your choice: ")
	input := readKey()

	switch string(input) {
	case "1":
		config := settings.Load()
		data := []byte(`{"pass":"` + config.UserPass + `",` + `"user":"` + config.UserName + `"` + `}`)

		req, err := http.NewRequest(http.MethodPost, "http://"+config.ConnectIP+":"+strconv.Itoa(config.ClientPort)+"/authorize", bytes.NewBuffer(data))
		if err != nil {
			panic(err)
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			//panic(err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
			}
		}(resp.Body)

		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		// TODO: Dodaj http request z tokenem wygenerowanym powyżej poprzez autoryzacje
		// TODO: Dodaj pętle, możliwość exploracji oraz pobierania plików
		// TODO: Dodaj wielowątkową opcje TCP do pobierania danych
		// TODO: Dodaj weryfikacje pobranych danych

		token := string(bodyBytes)

		fmt.Println("[client] Autoryzacja ukończona pomyślne") //\n[>>" + token + "<<]")

		// Test exploracji /
		req, err = http.NewRequest(http.MethodGet, "http://"+config.ConnectIP+":"+strconv.Itoa(config.ClientPort)+"/explore", bytes.NewBuffer(data))
		if req != nil {
			req.Header.Set("Token", token)
		}
		if err != nil {
			panic(err)
		}
		resp, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			//panic(err)
		}
		bodyBytes, err = io.ReadAll(resp.Body)
		fmt.Println("\033[31m" + string(bodyBytes) + "\u001B[0m")

		fmt.Println("[client] Pobieranie some.file.... \n some.file content:")
		req, err = http.NewRequest(http.MethodGet, "http://"+config.ConnectIP+":"+strconv.Itoa(config.ClientPort)+"/explore?path=/&file=some.file", bytes.NewBuffer(data))
		if req != nil {
			req.Header.Set("Token", token)
		}

		if err != nil {
			panic(err)
		}
		resp, err = client.Do(req)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
			//panic(err)
		}
		bodyBytes, err = io.ReadAll(resp.Body)
		fmt.Println("\033[31m" + string(bodyBytes) + "\u001B[0m")

		fmt.Println("\n[client] Zakonczone połączenie")

	case "2":
		listener.Host(-1)

	case "3":
		config := settings.Load()

		fmt.Println("")
		fmt.Printf("File loaded: %s\n", viper.ConfigFileUsed())
		fmt.Println("\tServer Config:")
		fmt.Printf("Server Port: %d\n", config.ServerPort)
		fmt.Printf("Server DB: %s\n", config.ServerDB)
		fmt.Printf("Shared: %s\n\n", config.Shared)

		fmt.Println("\tClient Config:")
		fmt.Printf("Connect IP: %s\n", config.ConnectIP)
		fmt.Printf("Connect Port: %d\n", config.ClientPort)
		fmt.Printf("Username: %s\n", config.UserName)
		fmt.Printf("Userpass: %s\n", config.UserPass)
		fmt.Printf("Downloads: %s\n", config.Downloads)

		fmt.Print("<< Press enter to continue\n")
		reader := bufio.NewReader(os.Stdin)
		_, _, _ = reader.ReadRune()

		Run()

	case "4":
		listener.ConfigDB()

	case "X":
		fmt.Println("Exiting...")
		os.Exit(0)

	default:
		fmt.Println("Invalid choice")
		fmt.Print("<< Press enter to continue\n")
		reader := bufio.NewReader(os.Stdin)
		_, _, _ = reader.ReadRune()
		Run()
	}

}
