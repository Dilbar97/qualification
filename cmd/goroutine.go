package cmd

import (
	"qualification/internal/usecase"

	"github.com/spf13/cobra"
)

var goroutineCmd = &cobra.Command{
	Use: "gor",

	Run: func(cmd *cobra.Command, args []string) {
		withChannel, err := cmd.Flags().GetBool("with-channel")
		if err != nil {
			return
		}

		withWG, err := cmd.Flags().GetBool("with-wg")
		if err != nil {
			return
		}

		withMutex, err := cmd.Flags().GetBool("with-mutex")
		if err != nil {
			return
		}

		usecase.StartGor(withChannel, withWG, withMutex)
	},
}

func init() {
	rootCmd.AddCommand(goroutineCmd)
	rootCmd.PersistentFlags().BoolP("with-channel", "c", false, "Run Goroutine with channel")
	rootCmd.PersistentFlags().BoolP("with-wg", "w", false, "Run Goroutine with wait group")
	rootCmd.PersistentFlags().BoolP("with-mutex", "m", false, "Run Goroutine with mutex")
}
