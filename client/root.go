package client

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
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

	fmt.Println("<<< SFE - Small File Exchanger >>>")
	fmt.Println("[1] Connect to Server")
	fmt.Println("[2] Host a server")
	fmt.Println("[3] Show config")
	fmt.Println("[4] Exit")
	//fmt.Println("Your choice: \"" + string(input) + "\"")
	input := readKey()

	switch string(input) {
	case "1":
		config := settings.Load()
		data := []byte(`{"pass":"` + config.ServerPass + `"}`)

		req, err := http.NewRequest(http.MethodPost, "http://"+config.ConnectIP+":"+strconv.Itoa(config.ClientPort)+"/authorize", bytes.NewBuffer(data))
		if err != nil {
			panic(err)
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {

			}
		}(resp.Body)

		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println(string(bodyBytes))
		fmt.Println("Zakonczone połączenie")

	case "2":
		listener.Host(8670)

	case "3":
		config := settings.Load()

		fmt.Println("")
		fmt.Printf("File loaded: %s\n", viper.ConfigFileUsed())
		fmt.Println("Server Config:")
		fmt.Printf("Server Port: %d\n", config.ServerPort)
		fmt.Printf("Server Password: %s\n\n", config.ServerPass)

		fmt.Println("Client Config:")
		fmt.Printf("Client Port: %d\n", config.ClientPort)
		fmt.Printf("Connect IP: %s\n", config.ConnectIP)

		fmt.Print("<< Press enter to continue\n")
		reader := bufio.NewReader(os.Stdin)
		_, _, _ = reader.ReadRune()

		Run()

	case "4":
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
