CREATE MIGRATION m1fjy4i5edtczsyuksemhlucpuz6nxrj2v5d7r6rd7xwyhaxlhzfpq
    ONTO m1sznfqu2wrzglz5x2v4mj2sfz2lpubvbxi47ev56rzfkksf6sr3fq
{
  ALTER TYPE default::Auditable {
      ALTER LINK meta {
          CREATE REWRITE
              UPDATE 
              USING (UPDATE
                  .meta
              SET {
                  modified := std::datetime_of_statement()
              });
      };
  };
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
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Class) THEN .parent ELSE .parent.class) ELSE <taxonomy::Taxon>{}));
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
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Class) THEN .parent ELSE .parent.class) ELSE <taxonomy::Taxon>{}));
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
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Family) THEN .parent ELSE .parent.family) ELSE <taxonomy::Taxon>{}));
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
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Family) THEN .parent ELSE .parent.family) ELSE <taxonomy::Taxon>{}));
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
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Genus) THEN .parent ELSE .parent.genus) ELSE <taxonomy::Taxon>{}));
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
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Genus) THEN .parent ELSE .parent.genus) ELSE <taxonomy::Taxon>{}));
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
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Kingdom) THEN .parent ELSE .parent.kingdom) ELSE <taxonomy::Taxon>{}));
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
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Kingdom) THEN .parent ELSE .parent.kingdom) ELSE <taxonomy::Taxon>{}));
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
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Order) THEN .parent ELSE .parent.order) ELSE <taxonomy::Taxon>{}));
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
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Order) THEN .parent ELSE .parent.order) ELSE <taxonomy::Taxon>{}));
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
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Phylum) THEN .parent ELSE .parent.phylum) ELSE <taxonomy::Taxon>{}));
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
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Phylum) THEN .parent ELSE .parent.phylum) ELSE <taxonomy::Taxon>{}));
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
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Species) THEN .parent ELSE .parent.species) ELSE <taxonomy::Taxon>{}));
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
              USING ((IF EXISTS (.parent) THEN (IF (.parent.rank = taxonomy::Rank.Species) THEN .parent ELSE .parent.species) ELSE <taxonomy::Taxon>{}));
      };
  };
};
