CREATE MIGRATION m1bboupmcsbvd2bmvgsqxnh6eukjluxyhdghutbovnog7h2ouj5iva
    ONTO m1azc6yiwh4btlicyyxmjrtz6qai7xkzm7cbbkecfvbenrgpwhaffq
{
  DROP FUNCTION location::search_sites(query: std::str, NAMED ONLY limit_n: OPTIONAL std::int32, NAMED ONLY threshold: std::float32);
  CREATE FUNCTION location::search_sites(query: std::str, NAMED ONLY threshold: std::float32 = 0.7) -> SET OF location::Site USING (SELECT
      location::Site {
          score := std::max({ext::pg_trgm::word_similarity(query, .name), ext::pg_trgm::word_similarity(query, .code), ext::pg_trgm::word_similarity(query, .locality)})
      }
  FILTER
      (.score > threshold)
  ORDER BY
      .score DESC
  );
  DROP FUNCTION location::sites_proximity(lat: std::float32, lon: std::float32, NAMED ONLY distance: std::int32, NAMED ONLY limit_n: OPTIONAL std::int32);
  CREATE FUNCTION location::sites_proximity(lat: std::float32, lon: std::float32, NAMED ONLY distance: std::int32) -> SET OF location::Site USING (SELECT
      location::Site
  FILTER
      (location::site_distance(location::Site, lat, lon) <= distance)
  );
};
