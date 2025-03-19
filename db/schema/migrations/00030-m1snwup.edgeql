CREATE MIGRATION m1snwup4t2r3namkkatu5a2ryxw33zyhnomrlrxcbjx4dy3f2vmfoq
    ONTO m1clh7odftbykcm4dqe7ythf6v3tuyg66c2fddku7gza2asor7kpza
{
  ALTER FUNCTION location::site_distance(site: location::Site, lat: std::float32, lon: std::float32) USING (std::assert_exists(ext::postgis::distance(ext::postgis::to_geography(location::WGS84_point(site.coordinates.latitude, site.coordinates.longitude)), ext::postgis::to_geography(location::WGS84_point(lat, lon)))));
};
