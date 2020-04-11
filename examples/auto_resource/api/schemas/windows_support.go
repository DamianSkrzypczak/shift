package schemas

import (
	"net/http"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

type fsWrapperForWindows struct {
	fs http.FileSystem
}

// Open replace windows-specific root for unix-like root path
func (fsw fsWrapperForWindows) Open(n string) (http.File, error) {
	return fsw.fs.Open(strings.Replace(n, "c:\\", "/", 1))
}

// newWindowsReferenceLoaderFileSystem decorates name with windows-compatible root path
// (in other hand gojsonschema throws "non canonical path" error)
// because filesystem passed to loader constructor is wrapped by fsWrapperForWindows,
// NewReferenceLoaderFileSystem receives valid windows path
// which then is converted back to unix path by fsWrapperForWindows.Open method
func newWindowsReferenceLoaderFileSystem(name string, fs http.FileSystem) gojsonschema.JSONLoader {
	return gojsonschema.NewReferenceLoaderFileSystem("file://c:/"+name, fsWrapperForWindows{fs})
}
