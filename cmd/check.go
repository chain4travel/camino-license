// Copyright (C) 2022-2024, Chain4Travel AG. All rights reserved.
// See the file LICENSE for licensing terms.

package cmd

import (
	"fmt"

	caminolicense "github.com/chain4travel/camino-license/pkg/camino-license"
	config "github.com/chain4travel/camino-license/pkg/config"
	"github.com/spf13/cobra"
)

// checkCmd represents the check command
var checkCmd = &cobra.Command{
	Use:   "check [FLAGS] FILES/DIRS",
	Short: "camino-license to check license headers",
	Long:  `camino-license to check license headers if they are compatible with the templates definded in the configuration file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configFile, _ := cmd.Flags().GetString("config")
		headersConfig, err := config.GetHeadersConfig(configFile)
		if err != nil {
			return err
		}
		h := caminolicense.CaminoLicenseHeader{Config: headersConfig}
		wrongFiles, err := h.CheckLicense(args)
		if err != nil {
			filesNumber := len(wrongFiles)
			if filesNumber == 1 {
				fmt.Println("1 file has wrong License Headers:")
			} else {
				fmt.Println(filesNumber, "files have wrong License Headers:")
			}
			for _, f := range wrongFiles {
				fmt.Println(f.File, "  - Reason:", f.Reason)
			}
			return err
		}
		fmt.Println("Check has finished successfully. All files have correct License Headers.")
		return nil
	},
}

// adding flags and check to camino-license command
func init() {
	checkCmd.Flags().StringP("config", "c", "config.yaml", "configuration yaml file path")
	rootCmd.AddCommand(checkCmd)
}
