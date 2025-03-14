package seeds

import (
	"context"
	"os"

	"github.com/geldata/gel-go/geltypes"
)

func SeedCountriesGeoJSON(db geltypes.Executor, path string) error {
	// Load JSON from file
	jsonBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return db.Execute(context.Background(),
		`#edgeql
			with module location,
			data := <json>$0,
			for item in json_array_unpack(data['features']) union (
				with country := (
					insert Country {
						name := <str>item['properties']['name'],
						continent := <str>item['properties']['region'],
						subcontinent := <str>item['properties']['sub-region'],
						code := (
							assert_exists(
								default::null_if_empty(<str>item['properties']['alpha-3']),
								message := "Country code is missing for " ++ <str>item['properties']['name']
							)
						)
					} unless conflict
				),
				select (
					if item['geometry'] != to_json("null") then (
						insert CountryBoundary {
							country := country,
							geometry := ext::postgis::geomfromgeojson(item['geometry'])
						} unless conflict
					) else (
						<CountryBoundary>{}
					)
				)
			)
		`, jsonBytes)
}
