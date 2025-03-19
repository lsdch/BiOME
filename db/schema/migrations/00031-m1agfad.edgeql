CREATE MIGRATION m1agfad6t2s4ghv5uhcobf4qzwnsmav5xyoljyk7h7pstp2w6nhd7a
    ONTO m1snwup4t2r3namkkatu5a2ryxw33zyhnomrlrxcbjx4dy3f2vmfoq
{
  ALTER FUNCTION location::site_distance(site: location::Site, lat: std::float32, lon: std::float32) USING (ext::postgis::distance(ext::postgis::to_geography(location::WGS84_point(site.coordinates.latitude, site.coordinates.longitude)), ext::postgis::to_geography(location::WGS84_point(lat, lon))));
  DROP FUNCTION location::sites_proximity(lat: std::float32, lon: std::float32, NAMED ONLY distance: std::int32);
};
