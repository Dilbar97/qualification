package cmd

import (
	"qualification/internal/handler"

	"github.com/spf13/cobra"

	_ "qualification/docs"
)

var swaggerCmd = &cobra.Command{
	Use: "rest",

	Run: func(cmd *cobra.Command, args []string) {
		withPgxLog, err := cmd.Flags().GetBool("log")
		if err != nil {
			return
		}

		withES, err := cmd.Flags().GetBool("es")
		if err != nil {
			return
		}

		a := handler.App{}
		a.Run(withPgxLog, withES)
	},
}

func init() {
	rootCmd.AddCommand(swaggerCmd)
	rootCmd.PersistentFlags().BoolP("log", "l", false, "Run mux handler with pgx log")
	rootCmd.PersistentFlags().BoolP("es", "e", false, "Run mux handler with elasticsearch")
}
