package settings

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort int
	ServerPass string
	ServerDB   string
	ClientPort int
	ConnectIP  string
	UserPass   string
	UserName   string
	Shared     string
	Downloads  string
}

func Load() Config {

	// Ustawianie domyślnych wartości konfiguracji
	viper.SetDefault("serverport", 7096)
	viper.SetDefault("serverpass", "password")
	viper.SetDefault("serverdb", "maindb")
	viper.SetDefault("clientport", 7096)
	viper.SetDefault("connectip", "localhost")
	viper.SetDefault("userpass", "password")
	viper.SetDefault("username", "user")
	viper.SetDefault("shared", "./share")
	viper.SetDefault("download", "./download")

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
		ServerPort: viper.GetInt("serverport"),
		ServerPass: viper.GetString("serverpass"),
		ServerDB:   viper.GetString("serverdb"),
		ClientPort: viper.GetInt("clientport"),
		ConnectIP:  viper.GetString("connectip"),
		UserPass:   viper.GetString("userpass"),
		UserName:   viper.GetString("username"),
		Shared:     viper.GetString("shared"),
		Downloads:  viper.GetString("download"),
	}

	return loaded
}

func PassVerify(pass string) bool {
	return viper.Get("serverpass").(string) == pass
}
