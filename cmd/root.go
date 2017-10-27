package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "tarkacoin",
	Short: "TarkaCoin is a sample blockchain implementaiton",
	Run: func(cmd *cobra.Command, args []string) {
	},
}
