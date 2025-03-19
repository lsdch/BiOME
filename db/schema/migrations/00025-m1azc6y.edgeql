CREATE MIGRATION m1azc6yiwh4btlicyyxmjrtz6qai7xkzm7cbbkecfvbenrgpwhaffq
    ONTO m17xsszprrpuykm3h2mdltiiqide7j5tsqaxohlfjqutw2ycui5b7a
{
  DROP FUNCTION location::search_sites_proximity(lat: std::float32, lon: std::float32, NAMED ONLY distance: std::int32, NAMED ONLY limit_n: OPTIONAL std::int32);
  DROP FUNCTION location::sites_proximity(lat: std::float32, lon: std::float32, distance: std::float32);
  CREATE FUNCTION location::sites_proximity(lat: std::float32, lon: std::float32, NAMED ONLY distance: std::int32, NAMED ONLY limit_n: OPTIONAL std::int32) -> SET OF location::Site USING (SELECT
      location::Site {
          distance := location::site_distance(location::Site, lat, lon)
      } FILTER
          (.distance <= distance)
      ORDER BY
          .distance ASC
  LIMIT
      limit_n
  );
};
