package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new content on the website",
	Run: func(cmd *cobra.Command, args []string) {
		isSite := isDirSite()
		if !isSite {
			fmt.Println("new: couldn't find config file 'config.toml'")
			return
		}

		if len(args) == 0 {
			fmt.Println("new: please provide a new name for the content")
			return
		}

		var createdFile *os.File
		defer createdFile.Close()

		path := args[0]
		if strings.Contains(path, "/") {
			parts := strings.Split(path, "/")

			// something -> contents/something.md
			// about/something -> contents/about/something.md
			filepath := fmt.Sprintf("./contents/%s", filepath.Join(parts[:len(parts)-1]...))
			err := os.MkdirAll(filepath, 0755)
			if err != nil {
				fmt.Println("new: error while creating new content")
				return
			}

			createdFile, err = os.Create(fmt.Sprintf("%s/%s.md", filepath, parts[len(parts)-1]))
			if err != nil {
				fmt.Println("new: error while creating new content file")
				return
			}

			return
		}

		createdFile, err := os.Create(fmt.Sprintf("./contents/%s.md", path))
		if err != nil {
			fmt.Println("new: error while creating new content file")
			return
		}

		title := convertToTitle(string(path[len(path)-1]))

		y, m, d := time.Now().Date()
		date := fmt.Sprintf("%s-%d-%d", m, d, y)
		contents := fmt.Sprintf("---\ntitle: \"%s\"\ndate: %s\n---", title, date)

		_, err = createdFile.Write([]byte(contents))
		if err != nil {
			fmt.Println("new: error while writing to file")
			return
		}
	},
}

func init() {
	newCmd.AddCommand(siteCmd)
}

func convertToTitle(slug string) string {
	var title []string
	str := "building-strings-in-golang"

	hasMultipleWords := strings.Contains(str, "-")
	if hasMultipleWords {
		parts := strings.Split(str, "-")
		for _, w := range parts {
			title = append(title, strings.ToUpper(string(w[0]))+w[1:])
		}
		return strings.Join(title, " ")
	}

	return strings.ToUpper(string(str[0])) + string(str[1:])
}

func isDirSite() bool {
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		return false
	}

	return true
}
