package routes

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/ekristen/pipeliner/pkg/box"
)

func (c *Routes) uiHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uri := strings.Replace(r.RequestURI, "/ui", "", -1)

		if uri == "" || uri == "/" {
			uri = "/index.html"
		}

		isHTML := false
		contentType := "text/plain"
		if strings.Contains(uri, ".css") {
			contentType = "text/css"
		} else if strings.Contains(uri, ".js") {
			contentType = "text/javascript"
		} else if strings.Contains(uri, ".png") {
			contentType = "image/png"
		} else if strings.Contains(uri, ".jpg") {
			contentType = "image/jpg"
		} else if strings.Contains(uri, ".html") {
			contentType = "text/html"
			isHTML = true
		}

		contents := box.Get(uri)
		if len(contents) == 0 {
			if !isHTML {
				contents = box.Get("/index.html")
				contentType = "text/html"
			} else {
				w.WriteHeader(404)
				return
			}
		}

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			var gr *gzip.Reader
			var err error
			gr, err = gzip.NewReader(bytes.NewBuffer(contents))
			if err != nil {
				c.log.WithError(err).Error("unable to decompress contents")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer gr.Close()
			contents, err = ioutil.ReadAll(gr)
			if err != nil {
				c.log.WithError(err).Error("unable to read uncompressed data")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		} else {
			w.Header().Add("content-encoding", "gzip")
		}

		w.Header().Add("content-length", fmt.Sprintf("%d", len(contents)))
		w.Header().Add("content-type", contentType)

		w.WriteHeader(200)
		w.Write(contents)
	})
}
