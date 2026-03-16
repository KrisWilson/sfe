package cmd

import (
	"sfe/listener"

	"github.com/spf13/cobra"
)

// TODO: exploreDir, downloadFile, downloadDir, uploadFile, uploadDir
// TODO: changePass from clientside

var exploreDirCmd = &cobra.Command{
	Use:   "explore",
	Short: "Pobiera informacje o zadanym folderze i zwraca informacje",
	Long:  "Pobiera informacje o zadanym folderze i zwraca informacje",
	//Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		listener.ConfigDB(2)
	},
}
