//go:generate statik -src=movies -p=movies -dest=. -include=*.json --ns movies -f
package schemas

import (
	"net/http"
	"runtime"
	"strings"

	"github.com/rakyll/statik/fs"
	"github.com/xeipuuv/gojsonschema"

	"app/api/schemas/movies"
)

func MustLoadMovieSchema(name string) gojsonschema.JSONLoader {
	schema, err := loadSchema(name, movies.Movies)
	if err != nil {
		panic(err)
	}

	return schema
}

const win = "windows"

// FSWrapper adapts statik file system (supporting unix-like paths) to windows paths
// by replacing "c:\" prefix with "/"
//

type FSWrapper struct { // Remove if you don't need windows support
	fs http.FileSystem
}

func (fsw FSWrapper) Open(name string) (http.File, error) { // Remove if you don't need windows support
	if runtime.GOOS == win {
		name = strings.Replace(name, "c:\\", "/", 1)
	}

	return fsw.fs.Open(name)
}

func loadSchema(name, namespace string) (gojsonschema.JSONLoader, error) {
	statikFS, err := fs.NewWithNamespace(namespace)
	if err != nil {
		return nil, err
	}

	path := "file:///" + name

	if runtime.GOOS == win { // Remove if you don't need windows support
		path = "file://c:/" + name
	}

	return gojsonschema.NewReferenceLoaderFileSystem(path, FSWrapper{statikFS}), nil
}
