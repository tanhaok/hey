/*
Copyright © 2024 Harold Hason
*/

package cmd

import (
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/sashabaranov/go-openai"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const GptKey = "CHAT_GPT_SECRET_KEY"

// gptCmd represents the gpt command
var gptCmd = &cobra.Command{
	Use:   "gpt",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKeySecret := os.Getenv(GptKey)
		inputText := strings.Join(args, " ")

		if apiKeySecret == "" {
			fmt.Println("API Secret doesn't exist")
			return
		}

		if inputText == "" {
			fmt.Println("Input text is empty. Existing...")
			return
		}

		chatWithGPT(inputText, apiKeySecret)

	},
}

func chatWithGPT(inputText string, apiKeyValue string) {
	client := openai.NewClient(apiKeyValue)
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5TurboInstruct,
		MaxTokens: 20,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: inputText,
			},
		},
		Stream: true,
	}

	stream, err := client.CreateChatCompletionStream(ctx, req)
	if err != nil {
		fmt.Printf("Completion stream error: %v\n", err)
		return
	}
	defer stream.Close()
	for {
		resp, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			fmt.Printf("Stream error: %v\n", err)
			return
		}

		color.Set(color.FgHiCyan)
		fmt.Println(resp.Choices[0].Delta.Content)
		color.Unset()

	}

}

func init() {
	rootCmd.AddCommand(gptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
