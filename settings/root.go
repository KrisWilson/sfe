package settings

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	serverPort int
	serverPass string
	clientPort int
	connectIP  string
}

func Load() Config {

	// Ustawianie domyślnych wartości konfiguracji
	viper.SetDefault("serverport", 7096)
	viper.SetDefault("serverpass", "password")
	viper.SetDefault("clientport", 7096)
	viper.SetDefault("connectip", "localhost")

	// Wczytanie ustawień z pliku settings.config
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig() // Wczytaj konfigurację z pliku
	if err != nil {
		fmt.Println("Error reading config file:", err)
		fmt.Println("Used default properties")
		err := viper.SafeWriteConfig()
		if err != nil {
			return Config{}
		}
	}

	loaded := Config{
		serverPort: viper.GetInt("serverport"),
		serverPass: viper.GetString("serverpass"),
		clientPort: viper.GetInt("clientport"),
		connectIP:  viper.GetString("connectip"),
	}

	fmt.Println("")
	fmt.Printf("File loaded: %s\n", viper.ConfigFileUsed())
	fmt.Println("Server Config:")
	fmt.Printf("Server Port: %d\n", loaded.serverPort)
	fmt.Printf("Server Password: %s\n\n", loaded.serverPass)

	fmt.Println("Client Config:")
	fmt.Printf("Client Port: %d\n", loaded.clientPort)
	fmt.Printf("Connect IP: %s\n", loaded.connectIP)

	fmt.Print("<< naciśnij enter aby kontynuuować\n")
	reader := bufio.NewReader(os.Stdin)
	_, _, _ = reader.ReadRune()

	return loaded
}
