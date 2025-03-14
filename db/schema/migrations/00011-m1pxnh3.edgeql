CREATE MIGRATION m1pxnh3jkudkehzohdcduh47ebzplqgng665u3zn3rhxtc6w5dkria
    ONTO m1md7emgbcxwcqirm5pmsrdyndreyz2q5mrt7wrkw47e6awahpikya
{
  CREATE FUNCTION location::WGS84_point(lon: std::float64, lat: std::float64) ->  ext::postgis::geometry USING (ext::postgis::point(lon, lat, 4326));
  ALTER FUNCTION location::position_to_country(lat: std::float32, lon: std::float32) USING (std::assert_single(((SELECT
      location::CountryBoundary
  FILTER
      ext::postgis::contains(location::CountryBoundary.geometry, location::WGS84_point(<std::float64>lon, <std::float64>lat))
  )).country, message := ((('More than one country found for the given position: ' ++ <std::str>lat) ++ ', ') ++ <std::str>lon)));
};
