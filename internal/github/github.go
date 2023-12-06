package github

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	gh "github.com/google/go-github/v57/github"
)

func GetInputFile(day int, useLocal bool) (io.ReadCloser, error) {
	if useLocal {
		dayDir := fmt.Sprintf("day%02d", day)
		path, err := filepath.Abs(filepath.Join("..", "internal", "inputs", dayDir, "input.txt"))
		if err != nil {
			return io.NopCloser(nil), fmt.Errorf("filepath.Abs: %w", err)
		}
		f, err := os.Open(path)
		if err != nil {
			return f, fmt.Errorf("os.Open(): %w", err)
		}
		return f, nil
	}
	token := os.Getenv("GITHUB_TOKEN")
	client := gh.NewClient(nil).WithAuthToken(token)
	input, resp, err := client.Repositories.DownloadContents(context.Background(),
		"harveysanders",
		"advent-of-code-inputs",
		"2023/day05/input.txt",
		&gh.RepositoryContentGetOptions{},
	)
	if err != nil {
		return input, fmt.Errorf("github.DownloadContents(): %w", err)
	}

	if resp.StatusCode >= 300 {
		log.Fatalf("%s", resp.Status)
		return input, fmt.Errorf("Github API HTTP: %s", resp.Status)
	}

	return input, nil
}
