package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var listAccountsCmd = &cobra.Command{
	Use:   "list-accounts",
	Short: "List Starling Bank accounts",
	Run: func(cmd *cobra.Command, args []string) {
		tokenFile := ".starling_token"
		tokenBytes, err := os.ReadFile(tokenFile)
		if err != nil {
			fmt.Println("Access token not found. Please run 'go-starling login' first.")
			os.Exit(1)
		}
		token := strings.TrimSpace(string(tokenBytes))

		client := &http.Client{}
		req, err := http.NewRequest("GET", "https://api.starlingbank.com/api/v2/accounts", nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			os.Exit(1)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			os.Exit(1)
		}

		if resp.StatusCode != 200 {
			fmt.Printf("API error: %s\n%s\n", resp.Status, string(body))
			os.Exit(1)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Error parsing JSON:", err)
			os.Exit(1)
		}

		jsonOut, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(jsonOut))
	},
}

var listTransactionsCmd = &cobra.Command{
	Use:   "list-transactions",
	Short: "List transactions for a Starling Bank account",
	Run: func(cmd *cobra.Command, args []string) {
		tokenFile := ".starling_token"
		tokenBytes, err := os.ReadFile(tokenFile)
		if err != nil {
			fmt.Println("Access token not found. Please run 'go-starling login' first.")
			os.Exit(1)
		}
		token := strings.TrimSpace(string(tokenBytes))

		accountUid, _ := cmd.Flags().GetString("account-uid")
		categoryUid, _ := cmd.Flags().GetString("category-uid")
		changesSince, _ := cmd.Flags().GetString("changes-since")

		if accountUid == "" || categoryUid == "" {
			fmt.Println("Both --account-uid and --category-uid flags are required")
			os.Exit(1)
		}

		client := &http.Client{}
		url := fmt.Sprintf("https://api.starlingbank.com/api/v2/feed/account/%s/category/%s?changesSince=%s", accountUid, categoryUid, changesSince)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println("Error creating request:", err)
			os.Exit(1)
		}
		req.Header.Set("Authorization", "Bearer "+token)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error making request:", err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response:", err)
			os.Exit(1)
		}

		if resp.StatusCode != 200 {
			fmt.Printf("API error: %s\n%s\n", resp.Status, string(body))
			os.Exit(1)
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Error parsing JSON:", err)
			os.Exit(1)
		}

		jsonOut, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(jsonOut))
	},
}

// ...existing code...

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
	// Default to 30 days ago for changesSince parameter
	defaultChangesSince := time.Now().UTC().AddDate(0, 0, -30).Format(time.RFC3339)
	
	listTransactionsCmd.Flags().String("account-uid", "", "Account UID (required)")
	listTransactionsCmd.Flags().String("category-uid", "", "Category UID (required)")
	listTransactionsCmd.Flags().String("changes-since", defaultChangesSince, "Timestamp to get transactions since (RFC3339 format)")
	
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(listAccountsCmd)
	rootCmd.AddCommand(listTransactionsCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
