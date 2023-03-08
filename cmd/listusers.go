/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// listusersCmd represents the listusers command
var listusersCmd = &cobra.Command{
	Use:   "listusers",
	Short: "List all users in GitLab SaaS",
	Long: `This command lists all users part of your GitLab SaaS instance.
`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		if url == "" {
			fmt.Println("GitLab instance URL is required. Use the --url flag to specify it.")
			os.Exit(1)
		}

		client := &http.Client{}
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v4/users", url), nil)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		req.Header.Set("PRIVATE-TOKEN", "<<YOUR_GITLAB_PERSONAL_ACCESS_TOKEN>>")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			fmt.Printf("Error: %s\n", resp.Status)
			os.Exit(1)
		}

		var users []struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
			Name     string `json:"name"`
			Email    string `json:"email"`
		}
		err = json.NewDecoder(resp.Body).Decode(&users)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		for _, u := range users {
			fmt.Printf("ID: %d, Username: %s, Name: %s, Email: %s\n", u.ID, u.Username, u.Name, u.Email)
		}
	},
}

func init() {
	rootCmd.AddCommand(listusersCmd)

	listusersCmd.Flags().String("url", "", "GitLab instance URL (required)")
	listusersCmd.MarkFlagRequired("url")
}
