CREATE MIGRATION m1diueiwkx6o5efowm6zzo3tjvcxprgolyvygiz6g4wi7ldvyuiooq
    ONTO m1o3detu65vtvhrjkrpmj46lnfc5flcqqtd334w6a3zliibiv6zdyq
{
  CREATE FUNCTION location::site_distance(site: location::Site, lat: std::float32, lon: std::float32) ->  std::float64 USING (ext::postgis::distance(ext::postgis::to_geography(location::WGS84_point(site.coordinates.latitude, site.coordinates.longitude)), ext::postgis::to_geography(location::WGS84_point(lat, lon))));
  ALTER FUNCTION location::sites_proximity(lat: std::float32, lon: std::float32, distance: std::float32) USING (SELECT
      location::Site
  FILTER
      (location::site_distance(location::Site, lat, lon) <= distance)
  );
};
