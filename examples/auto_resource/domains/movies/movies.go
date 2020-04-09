package movies

import (
	"fmt"
	"strconv"
	"time"

	"github.com/DamianSkrzypczak/shift"
	"github.com/DamianSkrzypczak/shift/autoapi"

	"app/api/schemas"
)

var movies = []Movie{}

type Movie struct {
	Title    string    `json:"title"`
	Genre    string    `json:"genre"`
	Released *JSONTime `json:"released"`
}

type JSONTime time.Time

func (t *JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", time.Time(*t).Format("2006-01-02"))), nil
}

func (t *JSONTime) UnmarshalJSON(data []byte) error {
	date, err := time.Parse("\"2006-01-02\"", string(data))

	if err != nil {
		return err
	}

	*t = JSONTime(date)

	return nil
}

type MovieCreationResponse struct {
	Link  string `json:"link"`
	Movie Movie  `json:"movie"`
}

func parseMovieID(id string) (int, error) {
	ID, err := strconv.Atoi(id)
	if err != nil {
		return 0, fmt.Errorf("no movie with index \"%s\"", id)
	}

	if ID < 0 || ID >= len(movies) {
		return 0, fmt.Errorf("no movie with index \"%d\"", ID)
	}

	return ID, nil
}

func Initialize(d *shift.Domain) {
	api := autoapi.NewResourceAPI(d)
	api.List(listMovies)
	api.Create(schemas.MustLoadMovieSchema("/create_request.json"), newCreateMovieHandler(api))
	api.Read(readMovie)
	api.Update(schemas.MustLoadMovieSchema("/update_request.json"), updateMovie)
	api.Replace(schemas.MustLoadMovieSchema("/replace_request.json"), replaceMovie)
	api.Delete(removeMovie)
}

func listMovies(params shift.QueryParameters) (interface{}, error) {
	return movies, nil
}

func newCreateMovieHandler(api *autoapi.ResourceAPI) func(deserializer autoapi.Deserializer, params shift.QueryParameters) (interface{}, error) {
	return func(deserializer autoapi.Deserializer, params shift.QueryParameters) (interface{}, error) {
		m := Movie{}
		if err := deserializer(&m); err != nil {
			return nil, err
		}

		movies = append(movies, m)
		movieResponse := MovieCreationResponse{api.ResourceURL(strconv.Itoa(len(movies) - 1)), m}

		return movieResponse, nil
	}
}

func readMovie(id string, params shift.QueryParameters) (interface{}, error) {
	index, err := parseMovieID(id)
	if err != nil {
		return nil, err
	}

	return movies[index], nil
}

func updateMovie(deserializer autoapi.Deserializer, id string, params shift.QueryParameters) (interface{}, error) {
	index, err := parseMovieID(id)
	if err != nil {
		return nil, err
	}

	m := movies[index]
	if err := deserializer(&m); err != nil {
		return nil, err
	}

	return nil, nil
}

func replaceMovie(deserializer autoapi.Deserializer, id string, params shift.QueryParameters) (interface{}, error) {
	index, err := parseMovieID(id)
	if err != nil {
		return nil, err
	}

	m := Movie{}
	if err := deserializer(&m); err != nil {
		return nil, err
	}

	movies[index] = m

	return nil, nil
}

func removeMovie(id string, params shift.QueryParameters) (interface{}, error) {
	index, err := parseMovieID(id)
	if err != nil {
		return nil, err
	}

	movies = append(movies[:index], movies[index+1:]...)

	return nil, nil
}
