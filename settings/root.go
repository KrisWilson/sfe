package settings

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ServerPort int
	ServerPass string
	ClientPort int
	ConnectIP  string
}

func Load() Config {

	// Ustawianie domyślnych wartości konfiguracji
	viper.SetDefault("serverport", 7096)
	viper.SetDefault("serverpass", "password")
	viper.SetDefault("clientport", 7096)
	viper.SetDefault("connectip", "localhost")
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
		ClientPort: viper.GetInt("clientport"),
		ConnectIP:  viper.GetString("connectip"),
	}

	return loaded
}
