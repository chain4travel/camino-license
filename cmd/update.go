package cmd

import (
	"fmt"

	caminolicense "github.com/chain4travel/camino-license/pkg/camino-license"
	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "camino-license to update license headers to the current year",
	Long:  `camino-license to update license headers to the current year if they are compatible with the templates definded in the configuration file`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configFile, _ := cmd.Flags().GetString("config")
		headersConfig, err := caminolicense.GetHeadersConfig(configFile)

		if err != nil {
			return err
		}

		updateErr := caminolicense.UpdateLicense(args, headersConfig)
		if updateErr != nil {
			return updateErr
		}
		fmt.Println("All License Headers have been updated successfully")
		return nil

	},
}

func init() {
	updateCmd.Flags().StringP("config", "c", "config.yaml", "configuration yaml file path")
	rootCmd.AddCommand(updateCmd)
}
