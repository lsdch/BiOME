package crossref

import (
	"darco/proto/models/people"
	"darco/proto/models/settings"

	"github.com/caltechlibrary/crossrefapi"
	"github.com/danielgtaylor/huma/v2"
	"github.com/edgedb/edgedb-go"
)

type CrossRefPerson crossrefapi.Person

func RetrieveDOI(db edgedb.Executor, doi string) (*crossrefapi.Works, error) {
	appName := settings.Instance().Name
	user, err := people.Current(db)
	if err != nil {
		return nil, err
	}

	client, err := crossrefapi.NewCrossRefClient(appName, user.Email)
	if err != nil {
		return nil, err
	}

	data, err := client.Works(doi)
	if data == nil && err == nil {
		return nil, huma.Error404NotFound("No match found")
	}

	return data, err
}
