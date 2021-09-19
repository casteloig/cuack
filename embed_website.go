package website

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	do "cuack/pkg/digitalocean"
)

var (
	//go:embed src/html
	content embed.FS

	//go:embed src/css
	css      embed.FS
	cssFS, _ = fs.Sub(css, "src/css")

	pages = map[string]string{
		"/": "src/html/inspect.html",
	}
)

var i do.Inspect

func CreateWebsite(inspected do.Inspect) error {

	i = inspected

	sm := http.NewServeMux()
	sm.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(cssFS))))
	x := http.FileServer(http.FS(content))
	sm.Handle("/src/html", x)
	sm.HandleFunc("/", rootHandler)

	server := http.Server{
		Handler: sm,
		Addr:    ":8080",
	}

	fmt.Println("You can connect to http://localhost:8080 to see the details")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	go serv(server)

	<-quit
	// gracefully stop server
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	fmt.Println("server stopped")
	return nil
}

func serv(server http.Server) error {
	err := server.ListenAndServe()
	if err != nil {
		return err
	}
	fmt.Println("server closed 100")

	return nil
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	page, ok := pages[r.URL.Path]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	tpl, err := template.ParseFS(content, page)
	if err != nil {
		log.Println("page not found in pages cache...")
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	if err := tpl.Execute(w, i); err != nil {
		log.Println(err)
		return
	}
}
