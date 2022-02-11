package frontend

import (
	"embed"
	"github.com/pierredavidbelanger/raftman/spi"
	"net/http"
	"net/url"
)

type uiFrontend struct {
	webFrontend
	api *apiFrontend
}

func newUIFrontend(e spi.LogEngine, frontendURL *url.URL) (*uiFrontend, error) {
	f := uiFrontend{}
	if err := initWebFrontend(e, frontendURL, &f.webFrontend); err != nil {
		return nil, err
	}
	f.api = &apiFrontend{}
	return &f, nil
}

//go:embed static/ui static/ui/index.html
var content embed.FS

func (f *uiFrontend) Start() error {
	_, b := f.e.GetBackend()
	f.api.b = b
	mux := http.NewServeMux()
	mux.HandleFunc(f.path+"api/stat", f.api.handleStat)
	mux.HandleFunc(f.path+"api/list", f.api.handleList)

	mux.Handle("/static/ui/", http.FileServer(http.FS(content)))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var con []uint8
		if r.RequestURI == "/" {
			con, _ = content.ReadFile("static/ui/index.html")
		} else {
			con, _ = content.ReadFile("static/ui" + r.RequestURI)
		}

		_, _ = w.Write(con)
	})

	return f.startHandler(mux)
}

func (f *uiFrontend) Close() error {
	return f.close()
}
