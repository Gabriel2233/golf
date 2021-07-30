package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	"github.com/spf13/cobra"
)

type MatterData struct {
	Title string
	Date  string
}

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new content on the website",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("new: please provide a name for the content, for example: golang/tutorial")
			return
		}

		if !hasConfig() {
			fmt.Println("new: couldn't find a config.toml here. Maybe you should create a new site with \"golf new site <SITE_NAME>\"")
			return
		}

		path := fmt.Sprintf("./contents/%s.md", args[0])
		err := ensureDir(path)
		if err != nil {
			fmt.Println("new: failed to create new content to the site: ", err)
			return
		}

		f, err := os.Create(path)
		if err != nil {
			fmt.Println("new: failed to create new content to the site: ", err)
			return
		}
		defer f.Close()

		filename := strings.TrimSuffix(filepath.Base(path), ".md")
		postTitle := formatTitle(filename)

		y, m, d := time.Now().Date()
		postDate := fmt.Sprintf("%s %d, %d", m, d, y)

		err = writeMatter(f, MatterData{postTitle, postDate})
		if err != nil {
			fmt.Println("new: failed to write default contents to file: ", err)
			return
		}

		fmt.Printf("%s was succesfully created\n", path)
	},
}

func init() {
	newCmd.AddCommand(siteCmd)
}

func formatTitle(slug string) string {
	var title strings.Builder
	slug = strings.ReplaceAll(slug, "-", " ")

	for i, r := range slug {
		if i == 0 {
			title.WriteRune(unicode.ToUpper(r))
			continue
		}

		if unicode.IsSpace(rune(slug[i-1])) {
			title.WriteRune(unicode.ToUpper(r))
			continue
		}

		title.WriteRune(r)
	}

	return title.String()
}

func writeMatter(w io.Writer, data MatterData) error {
	contents := fmt.Sprintf(`---
title: "%s"
date: "%s"
---`, data.Title, data.Date)

	_, err := w.Write([]byte(contents))
	if err != nil {
		return err
	}

	return nil
}

func ensureDir(fileName string) error {
	dirName := filepath.Dir(fileName)
	if _, err := os.Stat(dirName); err != nil {
		err = os.MkdirAll(dirName, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func hasConfig() bool {
	if _, err := os.Stat("config.toml"); os.IsNotExist(err) {
		return false
	}

	return true
}
