CREATE MIGRATION m1wsvo7dcmnfjhz3uaou6ud2ajg4v76ehueaq5jgw5ou5wmhqzgjaq
    ONTO m1bboupmcsbvd2bmvgsqxnh6eukjluxyhdghutbovnog7h2ouj5iva
{
  DROP FUNCTION location::search_sites(query: std::str, NAMED ONLY threshold: std::float32);
  CREATE FUNCTION location::search_sites(query: std::str, NAMED ONLY threshold: OPTIONAL std::float32 = 0.7) -> SET OF location::Site USING (SELECT
      location::Site {
          score := std::max({ext::pg_trgm::word_similarity(query, .name), ext::pg_trgm::word_similarity(query, .code), ext::pg_trgm::word_similarity(query, .locality)})
      }
  FILTER
      (.score > threshold)
  ORDER BY
      .score DESC
  );
  CREATE REQUIRED GLOBAL location::SITE_SEARCH_THRESHOLD -> std::float32 {
      SET default := 0.7;
  };
};
