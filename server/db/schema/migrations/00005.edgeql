CREATE MIGRATION m1yv6b46o6qyge2mdxsarylysutuzy3wngbnzwaissokz5dwj7gxjq
    ONTO m1emojwi2egcdl6utcg6nfsluhfzyifalybrx2uotnxa2a4o3lquaq
{
  ALTER TYPE taxonomy::Taxon {
      CREATE CONSTRAINT std::expression ON (NOT (std::contains(.name, ' '))) EXCEPT (((.rank = taxonomy::Rank.Species) OR (.rank = taxonomy::Rank.Subspecies)));
  };
};
