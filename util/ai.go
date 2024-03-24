package util

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func GenerateInitialChatResponse(db *sql.DB, rawData string) string {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEM_API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	basePrompt := os.Getenv("BASE_PROMPT")
	model := client.GenerativeModel("gemini-pro")
	prompt := genai.Text(basePrompt + rawData)
	iter := model.GenerateContentStream(ctx, prompt)

	var response strings.Builder
	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		content := getContent(resp)
		response.WriteString(content)
	}

	return response.String()
}

func getContent(resp *genai.GenerateContentResponse) string {
	var contentBuilder strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				content := fmt.Sprintf("%v", part)
				contentBuilder.WriteString(content)
			}
		}
	}
	return contentBuilder.String()
}
