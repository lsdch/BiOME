CREATE MIGRATION m1n7iewohrf5dazy4by7xh4tvzhqenfqcpirifg4vxaoc5faqdegva
    ONTO m14ndhr7yccrqi4rihpqmpunhf3emzcrd4drq5dqmamfefnsyplwfq
{
  ALTER TYPE taxonomy::Taxon {
      ALTER LINK class {
          ON TARGET DELETE ALLOW;
      };
      ALTER LINK family {
          ON TARGET DELETE ALLOW;
      };
      ALTER LINK genus {
          ON TARGET DELETE ALLOW;
      };
      ALTER LINK kingdom {
          ON TARGET DELETE ALLOW;
      };
      ALTER LINK order {
          ON TARGET DELETE ALLOW;
      };
      ALTER LINK phylum {
          ON TARGET DELETE ALLOW;
      };
      ALTER LINK species {
          ON TARGET DELETE ALLOW;
      };
  };
};
