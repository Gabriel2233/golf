package cmd

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
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

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launch a local webserver with all the contents",
	Run: func(cmd *cobra.Command, args []string) {
		posts := make(map[string]Post)

		// Parse all the posts
		// walk the contents dir
		// if current is file, parse just the header for the home page, put it on a struct and keep an array?
		// if is dir, ignore it

		// once i have all the data from the posts I can create a mux that will server this data in the template
		err := filepath.Walk("./contents", func(path string, info fs.FileInfo, err error) error {
			if !info.IsDir() {
				source, err := ioutil.ReadFile(path)
				path := strings.TrimSuffix(path[len("contents/"):], ".md")

				// contents/create-strings.md -> create-strings

				if err != nil {
					return err
				}

				post, err := parseMd(source)
				post.Path = path
				if err != nil {
					return err
				}

				posts[path] = post
			}

			return nil
		})
		if err != nil {
			fmt.Println("serve: error while triying to walk contents tree")
			return
		}

		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			html := `<html><body>
            <h1>My Site</h1>
            <section>
                {{range .}}
                    <div>
                        <a href="/posts/{{.Path}}"><h2>{{.Title}}</h2></a>
                        <p>{{.Date}}</p>
                    </div>
                {{end}}
            </section>
            </html></body>`

			tpl := template.Must(template.New("").Parse(html))

			tpl.Execute(w, posts)
		})
		mux.HandleFunc("/posts/", func(w http.ResponseWriter, r *http.Request) {
			postPath := string(r.URL.Path[len("/posts/"):])

			post := posts[postPath]

			html := `<html><body>
			<h1>{{.Title}}</h1>
			<section>
			<time>{{.Date}}</time>
            <div id="content">{{.Body}}</div>
			</section>

            <script>
                document.getElementById("content").innerHTML = {{.Body}};
            </script>
			</html></body>`

			tpl := template.Must(template.New("").Parse(html))

			tpl.Execute(w, post)
		})

		log.Fatal(http.ListenAndServe(":1414", mux))
	},
}

func parseMd(source []byte) (Post, error) {
	markdown := goldmark.New(
		goldmark.WithExtensions(
			meta.Meta,
		),
	)

	var buf bytes.Buffer
	context := parser.NewContext()
	if err := markdown.Convert(source, &buf, parser.WithContext(context)); err != nil {
		return Post{}, err
	}

	metaData := meta.Get(context)
	title := metaData["title"].(string)
	date := metaData["date"].(string)

	return Post{Title: title, Date: date, Body: buf.String()}, nil
}
