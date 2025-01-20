CREATE MIGRATION m1xzlkx4c2spo237cz7ypx4rmldvfcuqaisldfwcxy5i52z3z6b66a
    ONTO m1d2lb2mgvfwvlhge5fxwv3l7ik2podkkplwius3nmnuker6owo7ca
{
  CREATE TYPE datasets::AbstractDataset EXTENDING default::Auditable {
      CREATE REQUIRED MULTI LINK maintainers: people::Person;
      CREATE PROPERTY description: std::str;
      CREATE REQUIRED PROPERTY label: std::str {
          CREATE CONSTRAINT std::max_len_value(40);
          CREATE CONSTRAINT std::min_len_value(4);
      };
      CREATE REQUIRED PROPERTY slug: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  ALTER TYPE datasets::Dataset {
      DROP INDEX ON (.slug);
      DROP EXTENDING default::Auditable;
      EXTENDING datasets::AbstractDataset LAST;
  };
  ALTER TYPE datasets::Dataset {
      ALTER LINK maintainers {
          RESET CARDINALITY;
          RESET OPTIONALITY;
          DROP OWNED;
          RESET TYPE;
      };
      ALTER PROPERTY description {
          DROP OWNED;
          RESET TYPE;
      };
      ALTER PROPERTY label {
          ALTER CONSTRAINT std::max_len_value(40) {
              DROP OWNED;
          };
          RESET OPTIONALITY;
          ALTER CONSTRAINT std::min_len_value(4) {
              DROP OWNED;
          };
          DROP OWNED;
          RESET TYPE;
      };
      ALTER PROPERTY slug {
          ALTER CONSTRAINT std::exclusive {
              DROP OWNED;
          };
          RESET OPTIONALITY;
          DROP OWNED;
          RESET TYPE;
      };
  };
};
