package spa

import "net/http"

// FileServer represents the interface for serving static files from a filesystem
type FileServer interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}

// FileSystem represents the interface for for opening a single file
type FileSystem interface {
	Open(string) (http.File, error)
}

// NewFileServer creates a new FileServer from the FileSystem
func NewFileServer(fs FileSystem) (FileServer, error) {
	return http.FileServer(fs), nil
}
