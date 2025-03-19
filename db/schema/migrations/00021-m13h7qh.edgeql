CREATE MIGRATION m13h7qh3d4kosqbpk5qof5uaagmamrpsjdenigpieua7hkm6la6h3q
    ONTO m1diueiwkx6o5efowm6zzo3tjvcxprgolyvygiz6g4wi7ldvyuiooq
{
  CREATE FUNCTION location::search_sites(query: std::str, NAMED ONLY limit_n: std::int32 = 10, NAMED ONLY threshold: std::float32 = 0.7) -> SET OF location::Site USING (SELECT
      location::Site {
          score := std::max({ext::pg_trgm::word_similarity(query, .name), ext::pg_trgm::word_similarity(query, .code), ext::pg_trgm::word_similarity(query, .locality)})
      } FILTER
          (.score > threshold)
      ORDER BY
          .score DESC
  LIMIT
      limit_n
  );
  CREATE FUNCTION location::search_sites_proximity(lat: std::float32, lon: std::float32, NAMED ONLY distance: std::int32, NAMED ONLY limit_n: std::int32 = 10) -> SET OF location::Site USING (SELECT
      location::Site {
          distance := location::site_distance(location::Site, lat, lon)
      } FILTER
          (.distance <= distance)
      ORDER BY
          .distance ASC
  LIMIT
      limit_n
  );
  ALTER TYPE location::Site {
      CREATE INDEX ON (.name);
  };
  ALTER TYPE location::Site {
      CREATE INDEX ON (.locality);
  };
  ALTER TYPE location::Site {
      DROP INDEX ON (.code);
  };
};
