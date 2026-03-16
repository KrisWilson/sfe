package settings

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	// server
	ServerPort int
	ServerPass string
	ServerDB   string
	SharedDir  string
	// client
	ClientPort  int
	ConnectIP   string
	UserPass    string
	UserName    string
	DownloadDir string
}

func Load() Config {

	// Ustawianie domyślnych wartości konfiguracji
	viper.SetDefault("serverport", 7096)
	viper.SetDefault("serverpass", "password")
	viper.SetDefault("serverdb", "maindb")
	viper.SetDefault("shared", "./share")
	// klienta defaultowe wartości parametrów konfiguracji
	viper.SetDefault("clientport", 7096)
	viper.SetDefault("connectip", "localhost")
	viper.SetDefault("userpass", "password")
	viper.SetDefault("username", "user")
	viper.SetDefault("download", "./download")

	// Wczytanie ustawień z pliku settings.config
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig() // Wczytaj konfigurację z pliku
	if err != nil {
		fmt.Println("[Settings] Error reading config file:\r", err)
		fmt.Println("[Settings] Used default properties\r")
		err := viper.SafeWriteConfig()
		if err != nil {
			return Config{}
		}
	}

	loaded := Config{
		// server config properties
		ServerPort: viper.GetInt("serverport"),
		ServerPass: viper.GetString("serverpass"),
		ServerDB:   viper.GetString("serverdb"),
		SharedDir:  viper.GetString("shared"),
		// client config properties
		ClientPort:  viper.GetInt("clientport"),
		ConnectIP:   viper.GetString("connectip"),
		UserPass:    viper.GetString("userpass"),
		UserName:    viper.GetString("username"),
		DownloadDir: viper.GetString("download"),
	}

	return loaded
}
