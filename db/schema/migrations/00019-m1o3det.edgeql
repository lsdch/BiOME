CREATE MIGRATION m1o3detu65vtvhrjkrpmj46lnfc5flcqqtd334w6a3zliibiv6zdyq
    ONTO m1gkphgalqcyoeudgwixbuybsqgqou73ud6kgvojvsy7y5tjpeokrq
{
  DROP FUNCTION location::insert_site(data: std::json);
  DROP FUNCTION location::site_coords_to_country(coords: tuple<precision: location::CoordinatesPrecision, latitude: std::float32, longitude: std::float32>);
  DROP FUNCTION location::position_to_country(lat: std::float32, lon: std::float32);
  DROP FUNCTION location::WGS84_point(lon: std::float64, lat: std::float64);
  CREATE FUNCTION location::WGS84_point(lat: std::float64, lon: std::float64) ->  ext::postgis::geometry USING (ext::postgis::point(lon, lat, 4326));
  CREATE FUNCTION location::position_to_country(lat: std::float32, lon: std::float32) -> OPTIONAL location::Country USING (std::assert_single(((SELECT
      location::CountryBoundary
  FILTER
      ext::postgis::contains(location::CountryBoundary.geometry, location::WGS84_point(<std::float64>lat, <std::float64>lon))
  )).country, message := ((('More than one country found for the given position: ' ++ <std::str>lat) ++ ', ') ++ <std::str>lon)));
  CREATE FUNCTION location::sites_proximity(lat: std::float32, lon: std::float32, distance: std::float32) -> SET OF location::Site USING (SELECT
      location::Site
  FILTER
      (ext::postgis::distance(ext::postgis::to_geography(location::WGS84_point(lat, lon)), ext::postgis::to_geography(location::WGS84_point(.coordinates.latitude, .coordinates.longitude))) <= distance)
  );
  CREATE FUNCTION location::site_coords_to_country(coords: tuple<precision: location::CoordinatesPrecision, latitude: std::float32, longitude: std::float32>) -> OPTIONAL location::Country USING (SELECT
      location::position_to_country(coords.latitude, coords.longitude)
  );
  CREATE FUNCTION location::insert_site(data: std::json) ->  location::Site USING (WITH
      coords := 
          location::coords_from_json((data)['coordinates'])
      ,
      country_code := 
          <std::str>std::json_get(data, 'country_code')
  INSERT
      location::Site
      {
          name := <std::str>(data)['name'],
          code := <std::str>(data)['code'],
          description := <std::str>std::json_get(data, 'description'),
          coordinates := coords,
          locality := <std::str>std::json_get(data, 'locality'),
          country := (SELECT
              (IF EXISTS (country_code) THEN location::find_country(country_code) ELSE location::site_coords_to_country(coords))
          ),
          altitude := <std::int32>std::json_get(data, 'altitude')
      });
};
