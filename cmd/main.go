package main

import (
	"fmt"
	"os"

	"RestCLI/pkg"
	"RestCLI/web"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// Version information (set during build)
var (
	Version    = "dev"
	CommitHash = "unknown"
	BuildDate  = "unknown"
)

var rootCmd = &cobra.Command{
	Use:     "resterx-cli",
	Short:   "RESTerX CLI - A powerful CLI tool for testing APIs",
	Version: Version,
	Run: func(cmd *cobra.Command, args []string) {
		startInteractiveMenu()
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("RESTerX CLI\n")
		fmt.Printf("Version:    %s\n", Version)
		fmt.Printf("Commit:     %s\n", CommitHash)
		fmt.Printf("Build Date: %s\n", BuildDate)
	},
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the web interface",
	Long:  "Start the web server to access RESTerX through a web browser",
	Run: func(cmd *cobra.Command, args []string) {
		port, _ := cmd.Flags().GetString("port")
		web.StartWebServer(port)
	},
}

func init() {
	webCmd.Flags().StringP("port", "p", "8080", "Port to run the web server on")
	rootCmd.AddCommand(webCmd)
	rootCmd.AddCommand(versionCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func startInteractiveMenu() {
	for {
		prompt := promptui.Select{
			Label: "Select an HTTP method",
			Items: []string{"GET", "POST", "PUT", "PATCH", "HEAD", "DELETE", "Exit"},
		}

		_, choice, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed: %v\n", err)
			os.Exit(1)
		}

		switch choice {
		case "GET":
			fmt.Println("Selected GET method")
			getURL := promptURL()
			if getURL == "" {
				fmt.Println("URL cannot be empty")
				break
			}
			pkg.HandleGetRequest(getURL)
		case "POST":
			fmt.Println("Selected POST method")
			getURL := promptURL()
			if getURL == "" {
				fmt.Println("URL cannot be empty")
				break
			}
			pkg.HandlePostRequest(getURL)
		case "PUT":
			fmt.Println("Selected PUT method")
			getURL := promptURL()
			if getURL == "" {
				fmt.Println("URL cannot be empty")
				break
			}
			pkg.HandlePutRequest(getURL)
		case "PATCH":
			fmt.Println("Selected PATCH method")
			getURL := promptURL()
			if getURL == "" {
				fmt.Println("URL cannot be empty")
				break
			}
			pkg.HandlePatchRequest(getURL)
		case "HEAD":
			fmt.Println("Selected HEAD method")
			getURL := promptURL()
			if getURL == "" {
				fmt.Println("URL cannot be empty")
				break
			}
			pkg.HandleHeadRequest(getURL)
		case "DELETE":
			fmt.Println("Selected DELETE method")
			getURL := promptURL()
			if getURL == "" {
				fmt.Println("URL cannot be empty")
				break
			}
			pkg.HandleDeleteRequest(getURL)
			// Call function to handle DELETE requests
		case "Exit":
			fmt.Println("Exiting...")
			os.Exit(0)
		}
	}
}

func promptURL() string {
	prompt := promptui.Prompt{
		Label: "Enter URL:",
	}

	url, err := prompt.Run()
	if err != nil {
		fmt.Printf("Prompt failed: %v\n", err)
		os.Exit(1)
	}

	return url
}
