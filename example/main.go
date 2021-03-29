package main

import (
	"github.com/eliquious/console"
)

var (
	// Used for flags.
	cfgFile     string
	userLicense string
)

func main() {

	shell := console.New("mercator")
	scope := console.NewScope("binance", "Utilities for accessing the Binance crypto exchange")
	cmd := &console.Command{
		Use:   "risk",
		Short: "risk calculates an investment risk",
		Run: func(env *console.Environment, cmd *console.Command, args []string) error {
			return nil
		},
	}
	cmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.cobra.yaml)")
	cmd.Flags().StringP("author", "a", "YOUR NAME", "author name for copyright attribution")
	cmd.Flags().StringVarP(&userLicense, "license", "l", "", "name of license for the project")
	cmd.Flags().Bool("viper", true, "use Viper for configuration")

	scope.AddCommand(cmd)
	scope.AddSubScope(console.NewScope("account", "Access account info"))
	shell.AddScope(scope)
	shell.Run()
}
