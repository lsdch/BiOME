CREATE MIGRATION m1k4gv56i7bcctxglgneyvzsioeh7idvh5crj3lvgmpnok3mqkpz2q
    ONTO m1agfad6t2s4ghv5uhcobf4qzwnsmav5xyoljyk7h7pstp2w6nhd7a
{
  CREATE FUNCTION location::coords_distance(coords: tuple<precision: location::CoordinatesPrecision, latitude: std::float32, longitude: std::float32>, lat: std::float32, lon: std::float32) ->  std::float64 USING (ext::postgis::distance(ext::postgis::to_geography(location::WGS84_point(coords.latitude, coords.longitude)), ext::postgis::to_geography(location::WGS84_point(lat, lon))));
  ALTER FUNCTION location::site_distance(site: location::Site, lat: std::float32, lon: std::float32) USING (location::coords_distance(site.coordinates, lat, lon));
};
