package assets

import (
	"embed"
	"io"
	"io/fs"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

//go:embed static
var fsStatic embed.FS

// Register registers the ui on the root path.
func Register(r *mux.Router) {

	var assetsStatic, _ = fs.Sub(fsStatic, "static")
	var static = http.StripPrefix("/static/", http.FileServer(http.FS(assetsStatic)))
	r.PathPrefix("/static/").Handler(static)
}

func serveFile(fs fs.FS, name, contentType string) http.HandlerFunc {
	file, err := fs.Open(name)
	if err != nil {
		log.Panic().Err(err).Msgf("could not find %s", file)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		log.Panic().Err(err).Msgf("could not read %s", file)
	}

	return func(writer http.ResponseWriter, reg *http.Request) {
		writer.Header().Set("Content-Type", contentType)
		_, _ = writer.Write(content)
	}
}
