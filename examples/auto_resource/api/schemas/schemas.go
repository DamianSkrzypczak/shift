//go:generate statik -src=movies -p=movies -dest=. -include=*.json --ns movies -f
package schemas

import (
	"runtime"

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

func loadSchema(name, namespace string) (gojsonschema.JSONLoader, error) {
	statikFS, err := fs.NewWithNamespace(namespace)
	if err != nil {
		return nil, err
	}

	var path string
	if runtime.GOOS == "windows" {
		path = "file://d:/" + name
	} else {
		path = "file:///" + name
	}

	return gojsonschema.NewReferenceLoaderFileSystem(path, statikFS), nil
}
