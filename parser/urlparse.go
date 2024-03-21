package parser

import (
	"database/sql"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
	"golang.org/x/net/html"
)

// ParseURL fetches the content from the given URL and extracts all text within the body tag.
func ParseURL(db *sql.DB, link string) (string, error) {
	resp, err := http.Get(link)
	if err != nil {
		return "", fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML: %w", err)
	}

	var f func(*html.Node)
	var result strings.Builder
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			extractTextFromNode(n, &result)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	rawData := result.String()

	_, err = db.Exec("INSERT INTO scraped_data (link, raw_data) VALUES ($1, $2)", link, rawData)
	if err != nil {
		return "", fmt.Errorf("failed to insert data into database: %w", err)
	}

	return rawData, nil
}

// extractTextFromNode recursively extracts all text within the given node and its descendants.
func extractTextFromNode(n *html.Node, result *strings.Builder) {
	if n.Type == html.TextNode {
		result.WriteString(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractTextFromNode(c, result)
	}
}
