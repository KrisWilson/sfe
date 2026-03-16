package cmd

import (
	"sfe/client"
	"strings"

	"github.com/spf13/cobra"
)

// TODO: downloadFile, downloadDir, uploadFile, uploadDir
// TODO: changePass from clientside

var exploreDirCmd = &cobra.Command{
	Use:   "explore",
	Short: "[{path}] Pobiera informacje o zadanym folderze i zwraca informacje",
	Long:  "[{path}] Pobiera informacje o zadanym folderze i zwraca informacje",
	Run: func(cmd *cobra.Command, args []string) {
		client.ConnectServer()

		client.ExploreDir(strings.Join(args, " "))
	},
}
