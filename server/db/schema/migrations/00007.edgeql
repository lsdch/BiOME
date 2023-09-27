CREATE MIGRATION m1io5xqag7ldqlzmg2o7oxjrlh7uwdzx6g5pso2l2kh3fuzpblxgba
    ONTO m1eb5xvjiif2cw45hdn7aynjcojpsdzydjebeckbaika3347ikkscq
{
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK class {
          CREATE REWRITE
              INSERT 
              USING ((.parent.class IF EXISTS (.parent) ELSE <taxonomy::Taxon>{}));
          CREATE REWRITE
              UPDATE 
              USING ((.parent.class IF EXISTS (.parent) ELSE <taxonomy::Taxon>{}));
      };
      ALTER LINK family {
          CREATE REWRITE
              INSERT 
              USING ((.parent.family IF EXISTS (.parent) ELSE <taxonomy::Taxon>{}));
          CREATE REWRITE
              UPDATE 
              USING ((.parent.family IF EXISTS (.parent) ELSE <taxonomy::Taxon>{}));
      };
      ALTER LINK genus {
          CREATE REWRITE
              INSERT 
              USING ((.parent.genus IF EXISTS (.parent) ELSE <taxonomy::Taxon>{}));
          CREATE REWRITE
              UPDATE 
              USING ((.parent.genus IF EXISTS (.parent) ELSE <taxonomy::Taxon>{}));
      };
  };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK kingdom {
          DROP REWRITE
              INSERT ;
          };
      };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK kingdom {
          CREATE REWRITE
              INSERT 
              USING ((.parent IF (.parent.rank = taxonomy::Rank.Kingdom) ELSE (.parent.kingdom IF EXISTS (.parent) ELSE <taxonomy::Taxon>{})));
      };
  };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK kingdom {
          DROP REWRITE
              UPDATE ;
          };
      };
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK kingdom {
          CREATE REWRITE
              UPDATE 
              USING ((.parent IF (.parent.rank = taxonomy::Rank.Kingdom) ELSE (.parent.kingdom IF EXISTS (.parent) ELSE <taxonomy::Taxon>{})));
      };
      ALTER LINK order {
          CREATE REWRITE
              INSERT 
              USING ((.parent.order IF EXISTS (.parent) ELSE <taxonomy::Taxon>{}));
          CREATE REWRITE
              UPDATE 
              USING ((.parent.order IF EXISTS (.parent) ELSE <taxonomy::Taxon>{}));
      };
      ALTER LINK phylum {
          CREATE REWRITE
              INSERT 
              USING ((.parent.phylum IF EXISTS (.parent) ELSE <taxonomy::Taxon>{}));
          CREATE REWRITE
              UPDATE 
              USING ((.parent.phylum IF EXISTS (.parent) ELSE <taxonomy::Taxon>{}));
      };
      ALTER LINK species {
          CREATE REWRITE
              INSERT 
              USING ((.parent.species IF EXISTS (.parent) ELSE <taxonomy::Taxon>{}));
          CREATE REWRITE
              UPDATE 
              USING ((.parent.species IF EXISTS (.parent) ELSE <taxonomy::Taxon>{}));
      };
  };
};
