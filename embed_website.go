package website

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"

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

func CreateWebsite(inspected do.Inspect) error {

	// quit := make(chan os.Signal, 1)

	// signal.Notify(quit, os.Interrupt)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
		if err := tpl.Execute(w, inspected); err != nil {
			log.Println(err)
			return
		}
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(cssFS))))

	http.FileServer(http.FS(content))
	fmt.Println("You can connect to http://localhost:8080 to see the details")

	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			fmt.Println("Closing server")
			os.Exit(1)
		}
	}()

	return nil
}

// func gracefullShutdown(server *http.Server, logger *log.Logger, quit <-chan os.Signal) {
// 	<-quit
// 	logger.Println("Server is shutting down...")

// 	if err := server.Shutdown(ctx); err != nil {
// 	  logger.Fatalf("Could not gracefully shutdown the server: %v\n", err)
// 	}
//   }
