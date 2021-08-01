package cmd

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"time"

	"github.com/Gabriel2233/golf/pkg/http"
	"github.com/Gabriel2233/golf/pkg/markdown"
	"github.com/spf13/cobra"
)

type Post struct {
	Title string
	Date  string
	Body  string
	Path  string
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launch a local webserver with all the contents",
	Run: func(cmd *cobra.Command, args []string) {
		paths, err := getFilepaths()
		if err != nil {
			fmt.Println("serve: couldn't get filepaths")
			return
		}

		now := time.Now()
		posts := markdown.GetPosts(paths)
		fmt.Printf("took %s to get %d posts\n", time.Since(now), len(posts))

		http.LaunchServer(posts)
	},
}

func getFilepaths() ([]string, error) {
	var filepaths []string
	err := filepath.WalkDir("./contents/", func(path string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			filepaths = append(filepaths, path)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return filepaths, nil
}
