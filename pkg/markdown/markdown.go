package markdown

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"strings"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

type Post struct {
	Title string
	Date  string
	Body  string
	Path  string
}

type PostMatter struct {
	Title string
	Date  string
	Path  string
}

// this package is responsible for:
// reading the front-matter of a md file given a path and returning it
// reading the front-matter, as well as the post body, and returning it

func parseMatter(path string) (PostMatter, error) {
	f, err := os.Open(path)
	if err != nil {
		return PostMatter{}, err
	}
	defer f.Close()

	var source strings.Builder
	scanner := bufio.NewScanner(f)
	sepCount := 0

	for scanner.Scan() {
		if scanner.Text() == "---" && sepCount > 0 {
			source.WriteString(scanner.Text() + "\n")
			break
		}

		if scanner.Text() == "---" {
			source.WriteString(scanner.Text() + "\n")
			sepCount += 1
			continue
		}

		source.WriteString(scanner.Text() + "\n")
	}

	var buf bytes.Buffer
	md := goldmark.New(goldmark.WithExtensions(meta.Meta))

	context := parser.NewContext()
	err = md.Convert([]byte(source.String()), &buf, parser.WithContext(context))
	if err != nil {
		return PostMatter{}, err
	}

	metadata := meta.Get(context)
	matter := PostMatter{
		Title: metadata["title"].(string),
		Date:  metadata["date"].(string),
		Path:  path,
	}

	return matter, nil
}

func parseFullPost(path string) (Post, error) {
	f, err := os.Open(path)
	if err != nil {
		return Post{}, err
	}
	defer f.Close()

	source, err := ioutil.ReadAll(f)
	if err != nil {
		return Post{}, err
	}

	var buf bytes.Buffer
	md := goldmark.New(goldmark.WithExtensions(meta.Meta))

	context := parser.NewContext()
	err = md.Convert(source, &buf, parser.WithContext(context))
	if err != nil {
		return Post{}, err
	}

	metadata := meta.Get(context)
	post := Post{
		Title: metadata["title"].(string),
		Date:  metadata["date"].(string),
		Body:  buf.String(),
		Path:  path,
	}

	return post, nil
}

func GetAllMatters(paths []string) (matters []PostMatter, err error) {
	for _, path := range paths {
		matter, err := parseMatter(path)
		if err != nil {
			return nil, err
		}
		matters = append(matters, matter)
	}

	return matters, nil
}

func GetAllPosts(paths []string) (map[string]Post, error) {
	ret := make(map[string]Post, len(paths))
	for _, path := range paths {
		post, err := parseFullPost(path)
		if err != nil {
			return nil, err
		}
		ret[path] = post
	}

	return ret, nil
}
