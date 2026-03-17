package cmd

import (
	"fmt"
	"sfe/client"
	"strings"

	"github.com/spf13/cobra"
)

// TODO: downloadFile, downloadDir, uploadFile, uploadDir
// TODO: changePass from clientside

var downloadFileCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a file from Sfe",
	Run: func(cmd *cobra.Command, args []string) {
		// eg ./sfe download path/to/file/foo.bar
		client.ConnectServer()
		fmt.Println("Downloading file from Sfe ...")

		path := strings.Split(strings.Join(args, ""), "/")
		fmt.Println(path)
	},
}

var exploreDirCmd = &cobra.Command{
	Use:   "explore",
	Short: "[{path}] Pobiera informacje o zadanym folderze i zwraca informacje",
	Long:  "[{path}] Pobiera informacje o zadanym folderze i zwraca informacje",
	Run: func(cmd *cobra.Command, args []string) {
		client.ConnectServer()

		client.ExploreDir(strings.Join(args, " "))
	},
}
