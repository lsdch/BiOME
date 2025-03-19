CREATE MIGRATION m1ipdn3sij5enh5y4l7vkgk4kh3qifxr7rwtvjonc2qthlqmjdy2wq
    ONTO m1wsvo7dcmnfjhz3uaou6ud2ajg4v76ehueaq5jgw5ou5wmhqzgjaq
{
  DROP FUNCTION location::search_sites(query: std::str, NAMED ONLY threshold: OPTIONAL std::float32);
  CREATE FUNCTION location::search_sites(query: std::str, NAMED ONLY threshold: OPTIONAL std::float32 = <float64>{}) -> SET OF location::Site USING (SELECT
      location::Site {
          score := std::max({ext::pg_trgm::word_similarity(query, .name), ext::pg_trgm::word_similarity(query, .code), ext::pg_trgm::word_similarity(query, .locality)})
      }
  FILTER
      (.score > (threshold ?? GLOBAL location::SITE_SEARCH_THRESHOLD))
  ORDER BY
      .score DESC
  );
};
