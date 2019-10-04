package spa

import (
	"crypto/sha256"
	"fmt"
)

// Spa represents a single page application
type Spa interface {
	Index() []byte

	ETag() string

	StaticFileServer() *StaticServer
}

// StaticServer is used for creating a static file server in Gin
type StaticServer struct {
	Location   string
	FileSystem FileSystem
}

// New returns a default implementation of a Spa
func New(index []byte, location string, fs FileSystem) (Spa, error) {
	eTag := fmt.Sprintf("id_%x", sha256.Sum256(index))

	app := &baseApp{
		indexPage: index,
		eTag:      eTag,
		staticServer: &StaticServer{
			Location:   location,
			FileSystem: fs,
		},
	}

	return app, nil
}

type baseApp struct {
	indexPage    []byte
	eTag         string
	staticServer *StaticServer
}

// Index returns the cached index page
func (a *baseApp) Index() []byte {
	return a.indexPage
}

// StaticFileServer returns static file server
func (a *baseApp) StaticFileServer() *StaticServer {
	return a.staticServer
}

// StaticFileServer returns static file server
func (a *baseApp) ETag() string {
	return a.eTag
}
