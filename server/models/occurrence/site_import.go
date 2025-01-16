package occurrence

import (
	"context"
	"darco/proto/models/location"
	"darco/proto/services/geoapify"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/edgedb/edgedb-go"
)

type SiteImportDataset []SiteInput

func (d SiteImportDataset) RequestBody() []geoapify.LatLongCoords {
	body := make([]geoapify.LatLongCoords, len(d))
	for i, v := range d {
		body[i] = geoapify.LatLongCoords{
			Lat: v.Coordinates.Latitude,
			Lon: v.Coordinates.Longitude,
		}
	}
	return body
}

func (d SiteImportDataset) GroupByPrecision() map[location.CoordinatesPrecision][]SiteInput {
	var groups = map[location.CoordinatesPrecision][]SiteInput{
		location.M100:    []SiteInput{},
		location.KM1:     []SiteInput{},
		location.KM10:    []SiteInput{},
		location.KM100:   []SiteInput{},
		location.Unknown: []SiteInput{},
	}

	for _, v := range d {
		groups[v.Coordinates.Precision] = append(groups[v.Coordinates.Precision], v)
	}
	return groups
}

func (d SiteImportDataset) FillPlaces(apiKey string) error {
	// logrus.Infof("Input: \n %+v\n\n", d)
	// logrus.Infof("Body: \n %+v", d.RequestBody())
	// logrus.Exit(0)
	client := geoapify.NewGeoapifyClient(apiKey)
	response, err := client.BatchReverseGeocode(d.RequestBody())
	if err != nil {
		return err
	}

	// logrus.Infof("Geoapify response: \n %+v", response)

	for i, v := range *response {
		d[i].CountryCode = strings.ToUpper(v.CountryCode)
		d[i].Locality.SetValue(fmt.Sprintf("%s, %s", v.City, v.State))
	}
	return nil
}

func (i SiteImportDataset) Save(e edgedb.Executor) (created []Site, err error) {
	data, _ := json.Marshal(i)
	err = e.Query(context.Background(),
		fmt.Sprintf(`#edgeql
			with data := <json>$0,
			select (
				for site in json_array_unpack(data) union (
					%s
				)
			) { *, country: { * }, meta: { * }, datasets: { * } }
		`,
			i[0].InsertQuery("site"),
		), &created, data)
	return
}
