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
	Short: "List all users in a GitLab group",
	Long: `This command lists all users part of a GitLab group.
`,
	Run: func(cmd *cobra.Command, args []string) {
		url, _ := cmd.Flags().GetString("url")
		if url == "" {
			url = promptForInput("Enter GitLab instance URL: ")
		}

		token, _ := cmd.Flags().GetString("token")
		if token == "" {
			token = promptForInput("Enter GitLab personal access token: ")
		}

		group, _ := cmd.Flags().GetString("group")
		if group == "" {
			group = promptForInput("Enter GitLab group name: ")
		}

		client := &http.Client{}
		page := 1
		perPage := 200
		users := []struct {
			ID       int    `json:"id"`
			Username string `json:"username"`
			Name     string `json:"name"`
		}{}

		for {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v4/groups/%s/members?per_page=%d&page=%d", url, group, perPage, page), nil)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			req.Header.Set("PRIVATE-TOKEN", token)

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

			var pageUsers []struct {
				ID       int    `json:"id"`
				Username string `json:"username"`
				Name     string `json:"name"`
			}
			err = json.NewDecoder(resp.Body).Decode(&pageUsers)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if len(pageUsers) == 0 {
				break
			}

			users = append(users, pageUsers...)
			page++
		}

		for _, u := range users {
			fmt.Printf("ID: %d, Username: %s, Name: %s\n", u.ID, u.Username, u.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(listusersCmd)

	listusersCmd.Flags().String("url", "", "GitLab instance URL")
	listusersCmd.Flags().String("token", "", "GitLab personal access token")
	listusersCmd.Flags().String("group", "", "GitLab group name")
	listusersCmd.MarkFlagRequired("token")
}

func promptForInput(message string) string {
	var input string
	fmt.Print(message)
	fmt.Scanln(&input)
	return input
}
