package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"slices"
	"strings"

	"golang.org/x/net/html"
	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/components"
	. "maragu.dev/gomponents/html"
	ghttp "maragu.dev/gomponents/http"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	http.Handle("GET /{$}", http.RedirectHandler("/pt", http.StatusPermanentRedirect))
	http.Handle("GET /pt", HomeHandler("pt"))
	http.Handle("GET /en", HomeHandler("en"))
	http.Handle("GET /pt/tags/{tags}", TagsHandler("pt"))
	http.Handle("GET /en/tags/{tags}", TagsHandler("en"))
	http.Handle("GET /pt/posts/{slug}", PostHandler("pt"))
	http.Handle("GET /en/posts/{slug}", PostHandler("en"))
	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	// not found
	http.Handle("GET /pt/", PageNotFoundHandler("pt"))
	http.Handle("GET /en/", PageNotFoundHandler("en"))
	http.Handle("GET /", PageNotFoundHandler("en"))

	panic(http.ListenAndServe(":3000", nil))
}

type Post struct {
	title string
	desc  string
	slug  string
	href  string
	date  string
	tags  []string
	html  string
}

func loadPost(lang string, slug string) (Post, error) {
	filename := "./posts/" + lang + "/" + slug + ".html"

	b, err := os.ReadFile(filename)
	if err != nil {
		log.Print(err)
		return Post{}, err
	}

	raw := string(b)

	doc, err := html.Parse(strings.NewReader(raw))
	if err != nil {
		log.Print(err)
		return Post{}, err
	}

	h1 := htmlInnerText(htmlFirstByTag(doc, "h1"))
	if h1 == "" {
		err = fmt.Errorf("h1 not found for %s", filename)
		log.Print(err)
		return Post{}, err
	}

	p := htmlInnerText(htmlFirstByTag(doc, "p"))
	if p == "" {
		err = fmt.Errorf("p not found for %s", filename)
		log.Print(err)
		return Post{}, err
	}

	tags := strings.Split(htmlInnerText(htmlFirstByTag(doc, "tags")), " ")

	date, _, _ := strings.Cut(slug, "_")
	return Post{
		title: h1,
		desc:  p + " (...)",
		slug:  slug,
		date:  date,
		href:  "/" + lang + "/posts/" + slug,
		tags:  tags,
		html:  raw,
	}, nil
}

func loadPosts(lang string) ([]Post, error) {
	allEntries, err := os.ReadDir("./posts/" + lang)
	if err != nil {
		log.Print(err)
		return nil, err
	}

	posts := make([]Post, 0, len(allEntries))

	for _, e := range allEntries {
		if path.Ext(e.Name()) != ".html" {
			continue
		}

		post, err := loadPost(lang, strings.TrimSuffix(e.Name(), path.Ext(e.Name())))
		if err != nil {
			log.Print(err)
			continue
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func TagList(lang string, tags []string) Node {
	if len(tags) == 0 {
		return nil
	}

	return Div(
		Class("tag-list"),
		Map(tags, func(tag string) Node {
			return A(
				Class("link"),
				Href("/"+lang+"/tags/"+tag),
				Text(":"+tag+":"),
			)
		}),
	)
}

func PostList(lang string, posts []Post, tags []string) Node {
	return Map(posts, func(p Post) Node {
		t := translator(lang)
		show := true
		for _, tag := range tags {
			if !slices.Contains(p.tags, tag) {
				show = false
			}
		}

		if !show {
			return P(Text(t("No post found...")))
		}

		return Div(
			H2(A(Class("link"), Href(p.href), Text(p.title))),
			P(Text(p.date)),
			TagList(lang, p.tags),
			P(Text(p.desc)),
		)
	})
}

func Navbar(lang string, url string) Node {
	return Nav(
		A(Href("/"), Text("Home")),
		A(Href("https://www.linkedin.com/in/igoracmelo/"), Target("_blank"), Text("LinkedIn")),
		A(Href("https://github.com/igorcafe"), Target("_blank"), Text("GitHub")),
		A(Href("https://github.com/igorcafe/blog"), Target("_blank"), Text("Repo")),
		LangButton(lang, url),
	)
}

func LangButton(currLang string, url string) Node {
	var text string
	var newLang string
	switch currLang {
	case "en":
		text = "pt-BR"
		newLang = "pt"
	case "pt":
		text = "en-US"
		newLang = "en"
	default:
		log.Panic("unexpected language: ", currLang)
	}

	href := strings.Replace(url, "/"+currLang, "/"+newLang, 1)
	return A(Href(href), Text(text))
}

func HomeHandler(lang string) http.Handler {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		t := translator(lang)
		posts, err := loadPosts(lang)

		return Page(lang, "Blog", r.URL.String(), Main(
			H1(Text(t("Welcome!"))),
			P(Text(t("Here I talk about programming, (good and bad) experiences, Go, Linux, Emacs, and anything else that interests me."))),
			PostList(lang, posts, nil),
		)), err
	})
}

func TagsHandler(lang string) http.Handler {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		t := translator(lang)
		tags := strings.Split(r.PathValue("tags"), ",")
		posts, err := loadPosts(lang)

		return Page(lang, "Tags: "+strings.Join(tags, ", "), r.URL.String(), Main(
			H1(Text(t("Results for filter by tags: ")+strings.Join(tags, ", "))),
			PostList(lang, posts, tags),
		)), err
	})
}

func Page(lang string, title string, url string, content Node) Node {
	var langAttr string
	if lang == "pt" {
		langAttr = "pt-BR"
	} else {
		langAttr = "en-US"
	}

	return HTML5(HTML5Props{
		Title:       title + " | Igor Café",
		Description: "",
		Language:    langAttr,
		Head: Group{
			Link(Rel("stylesheet"), Href("/static/styles.css")),
		},
		Body: Group{
			Navbar(lang, url),
			content,
		},
	})
}

func PostHandler(lang string) http.Handler {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		t := translator(lang)
		post, err := loadPost(lang, r.PathValue("slug"))
		if err != nil {
			log.Print("load post: ", err)
			w.WriteHeader(http.StatusNotFound)
			return Page(lang, t("Post not found"), r.URL.String(), Main(
				H1(Text(t("Post not found"))),
				P(Text(t("Can't find the post you are looking for..."))),
			)), nil
		}

		return Page(lang, post.title, r.URL.String(), Raw(post.html)), nil
	})
}

func PostNotFoundHandler(lang string) http.Handler {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		t := translator(lang)

		return Page(lang, t("Post not found"), r.URL.String(), Main(
			H1(Text(t("Post not found"))),
			P(Text(t("Can't find the post you are looking for..."))),
		)), nil
	})
}

func PageNotFoundHandler(lang string) http.Handler {
	return ghttp.Adapt(func(w http.ResponseWriter, r *http.Request) (Node, error) {
		t := translator(lang)
		return Page(lang, t("Page not found"), r.URL.String(), Main(
			H1(Text(t("Page not found"))),
			P(Text(t("Can't find the page you are looking for..."))),
		)), nil
	})
}

func htmlWalk(root *html.Node, walk func(*html.Node) bool) *html.Node {
	if cont := walk(root); !cont {
		return root
	}

	var c *html.Node
	for c = root.FirstChild; c != nil; c = c.NextSibling {
		if found := htmlWalk(c, walk); found != nil {
			return found
		}
	}

	return c
}

func htmlFirstByTag(root *html.Node, tag string) *html.Node {
	return htmlWalk(root, func(n *html.Node) bool {
		return !(n.Type == html.ElementNode && n.Data == tag)
	})
}

func htmlInnerText(n *html.Node) string {
	if n == nil {
		return ""
	}
	var b strings.Builder

	htmlWalk(n, func(c *html.Node) bool {
		if c.Type == html.TextNode {
			b.WriteString(c.Data)
		}
		return true
	})

	return strings.TrimSpace(b.String())
}
