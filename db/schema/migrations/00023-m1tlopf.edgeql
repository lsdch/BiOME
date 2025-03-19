CREATE MIGRATION m1tlopfxqm3wbvgbdbq7wsslvlwouatl42zihbdylb7wcocfflv6sa
    ONTO m1nft4qlsmeygxinqz3ldicymno2xcfg3g55ol5qcmeexue267wl6a
{
  ALTER FUNCTION location::search_sites_proximity(lat: std::float32, lon: std::float32, NAMED ONLY distance: std::int32, NAMED ONLY limit_n: OPTIONAL std::int32 = 10) USING (SELECT
      location::Site {
          distance := location::site_distance(location::Site, lat, lon)
      } {
          distance
      } FILTER
          (.distance <= distance)
      ORDER BY
          .distance ASC
  LIMIT
      limit_n
  );
  CREATE ABSTRACT TYPE location::SiteWithDistance EXTENDING location::Site {
      CREATE REQUIRED PROPERTY distance: std::float64 {
          CREATE ANNOTATION std::description := 'Distance in meters from a reference point.';
      };
  };
};
