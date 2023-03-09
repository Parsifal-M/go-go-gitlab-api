/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// getReposCmd represents the getRepos command
var getReposCmd = &cobra.Command{
	Use:   "getRepos",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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

		var duration time.Duration
		if cmd.Flags().Changed("active") {
			active, _ := cmd.Flags().GetString("active")
			var err error
			duration, err = time.ParseDuration(active)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}

		client := &http.Client{}
		page := 1
		perPage := 50
		repos := []struct {
			Name              string `json:"name"`
			NameWithNamespace string `json:"name_with_namespace"`
			WebURL            string `json:"web_url"`
			LastActivityAt    string `json:"last_activity_at"`
		}{}

		for {
			req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v4/groups/%s/projects?per_page=%d&page=%d&include_subgroups=true", url, group, perPage, page), nil)
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

			var pageRepos []struct {
				Name              string `json:"name"`
				NameWithNamespace string `json:"name_with_namespace"`
				WebURL            string `json:"web_url"`
				LastActivityAt    string `json:"last_activity_at"`
			}

			if err := json.NewDecoder(resp.Body).Decode(&pageRepos); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			if len(pageRepos) == 0 {
				break
			}

			for _, r := range pageRepos {
				lastActivityTime, err := time.Parse(time.RFC3339, r.LastActivityAt)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				if cmd.Flags().Changed("active") {
					// Calculate the duration since the last activity time of the repository
					durationSinceLastActivity := time.Since(lastActivityTime)

					// Check if the duration since the last activity is less than the duration specified by the user
					if durationSinceLastActivity < duration {
						repos = append(repos, r)
					}
				} else {
					repos = append(repos, r)
				}
			}

			page++
		}

		for _, r := range repos {
			fmt.Printf("Name: %s\nRepo: %s\nWebURL: %s\nLastActivity: %s\n", r.Name, r.NameWithNamespace, r.WebURL, r.LastActivityAt)
		}
	},
}

func init() {
	rootCmd.AddCommand(getReposCmd)

	getReposCmd.Flags().StringP("url", "u", "", "GitLab instance URL")
	getReposCmd.Flags().StringP("token", "t", "", "GitLab personal access token")
	getReposCmd.Flags().StringP("group", "g", "", "GitLab group name")
	getReposCmd.Flags().StringP("active", "a", "", "Only return repos that have been active in the specified duration (e.g. 1h, 1d, 1w, 1m, 1y)")
}
