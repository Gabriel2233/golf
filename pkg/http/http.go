package http

import (
	"log"
	"net/http"
	"text/template"

	"github.com/Gabriel2233/golf/pkg/markdown"
)

// this package is responsible for launcing the webserver for the site
// i need a way to get the posts here, so i can render them in the template

func LaunchServer(posts map[string]markdown.Post) {
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

		var postsArr []markdown.Post

		for _, v := range posts {
			postsArr = append(postsArr, v)
		}

		tpl.Execute(w, postsArr)
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

	log.Print("server up")
	log.Fatal(http.ListenAndServe(":1414", mux))
}
