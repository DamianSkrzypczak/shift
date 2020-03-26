//go:generate statik -src=movies -p=movies -dest=. -include=*.json --ns movies -f
package schemas

import (
	"io/ioutil"

	"github.com/rakyll/statik/fs"

	"app/api/schemas/movies"
)

func MustLoadMovieSchema(name string) []byte {
	schema, err := loadSchema(name, movies.Movies)
	if err != nil {
		panic(err)
	}

	return schema
}

func loadSchema(name, namespace string) ([]byte, error) {
	statikFS, err := fs.NewWithNamespace(namespace)
	if err != nil {
		return nil, err
	}

	r, err := statikFS.Open(name)
	if err != nil {
		return nil, err
	}

	defer r.Close()

	contents, err := ioutil.ReadAll(r)

	if err != nil {
		return contents, err
	}

	return contents, nil
}
