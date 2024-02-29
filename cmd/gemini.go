/*
Copyright © 2024 HAROLD HASON
*/

package cmd

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const GeminiKey = "GEMINI_SECRET_KEY"

// geminiCmd represents the gemini command
var geminiCmd = &cobra.Command{
	Use:   "gemini",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,

	Run: func(cmd *cobra.Command, args []string) {
		chatWithGemini(cmd, args)
	},
}

func chatWithGemini(cmd *cobra.Command, args []string) {
	apiValue := os.Getenv(GeminiKey)

	if apiValue == "" {
		fmt.Println("API key for gemini does not exit.")
		os.Exit(1)
	}

	inputText := strings.Join(args, " ")
	if inputText == "" {
		fmt.Println("Input text is empty. Existing...")
		return
	} else {
		ctx := context.Background()
		client, err := genai.NewClient(ctx, option.WithAPIKey(apiValue))
		if err != nil {
			log.Fatal(err)
		}

		defer client.Close()

		model := client.GenerativeModel("gemini-pro")
		prompt := genai.Text(inputText)
		iter := model.GenerateContentStream(ctx, prompt)
		for {
			resp, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			color.Set(color.FgHiCyan)
			fmt.Println(resp.Candidates[0].Content.Parts[0])
			color.Unset()
		}

	}

}

func init() {
	rootCmd.AddCommand(geminiCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// geminiCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// geminiCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
