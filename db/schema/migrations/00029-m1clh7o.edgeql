CREATE MIGRATION m1clh7odftbykcm4dqe7ythf6v3tuyg66c2fddku7gza2asor7kpza
    ONTO m1ipdn3sij5enh5y4l7vkgk4kh3qifxr7rwtvjonc2qthlqmjdy2wq
{
  DROP FUNCTION location::search_sites(query: std::str, NAMED ONLY threshold: OPTIONAL std::float32);
  CREATE FUNCTION location::site_fuzzy_search_score(site: location::Site, query: std::str) ->  std::float32 USING (std::max({ext::pg_trgm::word_similarity(query, site.name), ext::pg_trgm::word_similarity(query, site.code), ext::pg_trgm::word_similarity(query, site.locality)}));
};
