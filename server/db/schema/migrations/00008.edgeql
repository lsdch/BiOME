CREATE MIGRATION m14ndhr7yccrqi4rihpqmpunhf3emzcrd4drq5dqmamfefnsyplwfq
    ONTO m1io5xqag7ldqlzmg2o7oxjrlh7uwdzx6g5pso2l2kh3fuzpblxgba
{
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK class {
          DROP REWRITE
              INSERT ;
          };
      };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK class {
          CREATE REWRITE
              INSERT 
              USING ((.parent IF (.parent.rank = taxonomy::Rank.Class) ELSE (.parent.class IF EXISTS (.parent) ELSE <taxonomy::Taxon>{})));
      };
  };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK class {
          DROP REWRITE
              UPDATE ;
          };
      };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK class {
          CREATE REWRITE
              UPDATE 
              USING ((.parent IF (.parent.rank = taxonomy::Rank.Class) ELSE (.parent.class IF EXISTS (.parent) ELSE <taxonomy::Taxon>{})));
      };
  };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK family {
          DROP REWRITE
              INSERT ;
          };
      };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK family {
          CREATE REWRITE
              INSERT 
              USING ((.parent IF (.parent.rank = taxonomy::Rank.Family) ELSE (.parent.family IF EXISTS (.parent) ELSE <taxonomy::Taxon>{})));
      };
  };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK family {
          DROP REWRITE
              UPDATE ;
          };
      };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK family {
          CREATE REWRITE
              UPDATE 
              USING ((.parent IF (.parent.rank = taxonomy::Rank.Family) ELSE (.parent.family IF EXISTS (.parent) ELSE <taxonomy::Taxon>{})));
      };
  };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK genus {
          DROP REWRITE
              INSERT ;
          };
      };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK genus {
          CREATE REWRITE
              INSERT 
              USING ((.parent IF (.parent.rank = taxonomy::Rank.Genus) ELSE (.parent.genus IF EXISTS (.parent) ELSE <taxonomy::Taxon>{})));
      };
  };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK genus {
          DROP REWRITE
              UPDATE ;
          };
      };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK genus {
          CREATE REWRITE
              UPDATE 
              USING ((.parent IF (.parent.rank = taxonomy::Rank.Genus) ELSE (.parent.genus IF EXISTS (.parent) ELSE <taxonomy::Taxon>{})));
      };
  };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK order {
          DROP REWRITE
              INSERT ;
          };
      };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK order {
          CREATE REWRITE
              INSERT 
              USING ((.parent IF (.parent.rank = taxonomy::Rank.Order) ELSE (.parent.order IF EXISTS (.parent) ELSE <taxonomy::Taxon>{})));
      };
  };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK order {
          DROP REWRITE
              UPDATE ;
          };
      };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK order {
          CREATE REWRITE
              UPDATE 
              USING ((.parent IF (.parent.rank = taxonomy::Rank.Order) ELSE (.parent.order IF EXISTS (.parent) ELSE <taxonomy::Taxon>{})));
      };
  };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK phylum {
          DROP REWRITE
              INSERT ;
          };
      };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK phylum {
          CREATE REWRITE
              INSERT 
              USING ((.parent IF (.parent.rank = taxonomy::Rank.Phylum) ELSE (.parent.phylum IF EXISTS (.parent) ELSE <taxonomy::Taxon>{})));
      };
  };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK phylum {
          DROP REWRITE
              UPDATE ;
          };
      };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK phylum {
          CREATE REWRITE
              UPDATE 
              USING ((.parent IF (.parent.rank = taxonomy::Rank.Phylum) ELSE (.parent.phylum IF EXISTS (.parent) ELSE <taxonomy::Taxon>{})));
      };
  };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK species {
          DROP REWRITE
              INSERT ;
          };
      };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK species {
          CREATE REWRITE
              INSERT 
              USING ((.parent IF (.parent.rank = taxonomy::Rank.Species) ELSE (.parent.species IF EXISTS (.parent) ELSE <taxonomy::Taxon>{})));
      };
  };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK species {
          DROP REWRITE
              UPDATE ;
          };
      };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK species {
          CREATE REWRITE
              UPDATE 
              USING ((.parent IF (.parent.rank = taxonomy::Rank.Species) ELSE (.parent.species IF EXISTS (.parent) ELSE <taxonomy::Taxon>{})));
      };
  };
};
