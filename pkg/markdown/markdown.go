package markdown

import (
	"bytes"
	"io"
	"os"

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

// this package is responsible for:
// reading the front-matter of a md file given a path and returning it
// reading the front-matter, as well as the post body, and returning it

func parsePost(path string) (Post, error) {
	f, err := os.Open(path)
	if err != nil {
		return Post{}, err
	}
	defer f.Close()

	source, err := io.ReadAll(f)
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

func GetPosts(paths []string) map[string]Post {
	ret := make(map[string]Post, len(paths))

	for _, path := range paths {
		post, err := parsePost(path)
		if err != nil {
			continue
		}

		ret[post.Path] = post
	}
	return ret
}
