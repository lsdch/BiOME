CREATE MIGRATION m1qhc5m73dsddd6uqldwbmkgfrdcyzkrwd6fvyjybvgjs2ndh5ieoq
    ONTO m1k4gv56i7bcctxglgneyvzsioeh7idvh5crj3lvgmpnok3mqkpz2q
{
  CREATE FUNCTION location::site_as_point(site: location::Site) ->  ext::postgis::geometry USING (location::WGS84_point(<std::float64>site.coordinates.latitude, <std::float64>site.coordinates.longitude));
};
