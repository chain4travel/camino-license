package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "camino-license COMMAND FLAGS",
	Short: "camino-license pkg to check and update license headers",
	Long:  `camino-license pkg to check and update license headers according to a given yaml configuration`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
