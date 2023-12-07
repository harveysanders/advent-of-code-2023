package github

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	gh "github.com/google/go-github/v57/github"
)

type File struct {
	io.Closer
	io.ReadSeeker
}

// GetInputFile fetches a file from the local disk if useLocal is true or from my private advent-of-code-inputs GitHub repo.
// The reader returned can be reset to the beginning.
//
// Ex:
//
//	r, _ := GetInputFile(1, true)
//	r.Seek(0, io.SeekStart)
func GetInputFile(day int, useLocal bool) (io.ReadSeekCloser, error) {
	dayDir := fmt.Sprintf("day%02d", day)
	if useLocal {
		path, err := filepath.Abs(filepath.Join("..", "internal", "inputs", dayDir, "input.txt"))
		if err != nil {
			return File{}, fmt.Errorf("filepath.Abs: %w", err)
		}
		f, err := os.Open(path)
		if err != nil {
			return f, fmt.Errorf("os.Open(): %w", err)
		}
		return f, nil
	}
	token := os.Getenv("GITHUB_TOKEN")
	client := gh.NewClient(nil).WithAuthToken(token)
	f, resp, err := client.Repositories.DownloadContents(context.Background(),
		"harveysanders",
		"advent-of-code-inputs",
		fmt.Sprintf("2023/%s/input.txt", dayDir),
		&gh.RepositoryContentGetOptions{},
	)
	if err != nil {
		return File{}, fmt.Errorf("github.DownloadContents(): %w", err)
	}

	if resp.StatusCode >= 300 {
		log.Fatalf("%s", resp.Status)
		return File{}, fmt.Errorf("Github API HTTP: %s", resp.Status)
	}

	// Create a ReadSeeker so the stream can be reset (read again).
	b, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	input := bytes.NewReader(b)

	return File{f, input}, nil
}
