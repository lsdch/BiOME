CREATE MIGRATION m1bfq5lgnnheudvkk5woblpwtebljmzytqglmkhbzvgrxesau7n4mq
    ONTO m1aabgqjxm74kmco3jgssfyalpgzi6mf4kqs2fogkccmqxqipiosrq
{
  CREATE FUNCTION location::position_to_country(lat: std::float32, lon: std::float32) -> OPTIONAL location::Country USING (std::assert_single(((SELECT
      location::CountryBoundary
  FILTER
      ext::postgis::contains(location::CountryBoundary.geometry, ext::postgis::point(<std::float64>lon, <std::float64>lat))
  )).country, message := ((('More than one country found for the given position: ' ++ <std::str>lat) ++ ', ') ++ <std::str>lon)));
};
