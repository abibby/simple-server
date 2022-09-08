package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

func main() {
	port := 3000
	dir := ""

	flag.StringVar(&dir, "dir", ".", "the directory to serve")
	flag.IntVar(&port, "port", 3000, "the port to serve on")

	flag.Parse()

	path, err := filepath.Abs(dir)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Serving %s at http://localhost:%d\n", path, port)

	http.ListenAndServe(fmt.Sprintf(":%d", port), FileServerDefault(os.DirFS(path), "", "index.html"))
}

func FileServerDefault(root fs.FS, basePath, fallbackPath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := path.Join(basePath, path.Clean(r.URL.Path)[1:])

		info, err := fs.Stat(root, p)
		if err != nil || info.IsDir() {
			b, err := fs.ReadFile(root, path.Join(basePath, fallbackPath))
			if err != nil {
				log.Print(err)
				return
			}

			w.Header().Add("Content-Type", "text/html")
			_, err = w.Write(b)
			if err != nil {
				log.Print(err)
				return
			}
			return
		}

		f, err := root.Open(p)
		if err != nil {
			log.Print(err)
			return
		}

		w.Header().Add("Content-Type", mime.TypeByExtension(path.Ext(p)))

		_, err = io.Copy(w, f)
		if err != nil {
			log.Print(err)
			return
		}
	})
}
