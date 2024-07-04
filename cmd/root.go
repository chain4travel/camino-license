// Copyright (C) 2022-2024, Chain4Travel AG. All rights reserved.
// See the file LICENSE for licensing terms.

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// camino-license command
var rootCmd = &cobra.Command{
	Use:   "camino-license COMMAND [FLAGS] FILES/DIRS",
	Short: "camino-license pkg to check license headers",
	Long:  `camino-license pkg to check license headers according to a given yaml configuration`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
