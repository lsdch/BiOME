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

func (d SiteImportDataset) FillPlaces(db edgedb.Executor, apiKey string) error {
	client := geoapify.NewGeoapifyClient(apiKey)
	response, err := client.BatchReverseGeocode(db, d.RequestBody())
	if err != nil {
		return err
	}

	for i, v := range *response {
		d[i].CountryCode = strings.ToUpper(v.CountryCode)
		switch d[i].Coordinates.Precision {
		case location.M100, location.KM1, location.KM10:
			locality := v.City
			if v.State != "" {
				locality = locality + ", " + v.State
			}
			d[i].Locality.SetValue(locality)
		case location.KM100:
			d[i].Locality.SetValue(v.State)
		case location.Unknown:
			// skip
			break
		}
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
