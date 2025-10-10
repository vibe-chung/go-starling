
package main

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "go-starling",
	Short: "Starling Bank CLI",
	Long:  `A CLI for interacting with Starling Bank APIs.`,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Save your Starling Bank API access token",
	Run: func(cmd *cobra.Command, args []string) {
		tokenFile := ".starling_token"
		if _, err := os.ReadFile(tokenFile); err == nil {
			fmt.Println("Access token already saved.")
			return
		}
		var token string
		fmt.Print("Enter your Starling Bank API access token: ")
		fmt.Scanln(&token)
		err := os.WriteFile(tokenFile, []byte(token), 0600)
		if err != nil {
			fmt.Println("Error saving token:", err)
			os.Exit(1)
		}
		fmt.Println("Access token saved successfully.")
	},
}

func main() {
	rootCmd.AddCommand(loginCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
