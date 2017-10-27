package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tarkalabs/tarkacoin/wallet"
)

var WalletCmd = &cobra.Command{
	Use: "wallet",
}
var walletName string

func init() {
	RootCmd.AddCommand(WalletCmd)
	newCommand := setupNewCommand()
	newCommand.PersistentFlags().StringVar(&walletName, "name", "n", "name of the wallet")
	WalletCmd.AddCommand(newCommand)
}

func setupNewCommand() *cobra.Command {
	return &cobra.Command{
		Use: "new",
		Run: func(cmd *cobra.Command, args []string) {
			err := wallet.SaveKey(walletName)
			if err != nil {
				panic(err)
			}
		},
	}
}
