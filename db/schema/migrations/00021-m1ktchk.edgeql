CREATE MIGRATION m1ktchkrmzvi33nnju4pes4laahsrlxpypklkacd5zzfquqb3vbrwa
    ONTO m12moatuy6u3hxui2xrsqhctqep2uywc7kfpwu2iaoasjd2repfpha
{
  CREATE MODULE references IF NOT EXISTS;
  CREATE MODULE traits IF NOT EXISTS;
  ALTER TYPE occurrence::Identification {
      ALTER LINK identified_by {
          RESET OPTIONALITY;
      };
  };
  ALTER TYPE occurrence::OccurrenceReport {
      CREATE PROPERTY content_description: std::str;
      ALTER PROPERTY quantity {
          DROP CONSTRAINT std::expression ON ((((.precision = occurrence::QuantityType.Exact) AND (.exact > 0)) OR (.precision != occurrence::QuantityType.Exact)));
          SET TYPE occurrence::QuantityType USING (<occurrence::QuantityType>.quantity.precision);
          DROP REWRITE
              INSERT ;
              DROP REWRITE
                  UPDATE ;
              };
          };
  ALTER TYPE occurrence::OccurrenceReport {
      ALTER PROPERTY specimen_available {
          RENAME TO voucher;
      };
  };
  ALTER TYPE reference::Article {
      DROP PROPERTY year;
  };
  ALTER TYPE reference::Article RENAME TO references::Article;
  ALTER TYPE references::Article {
      ALTER PROPERTY authors {
          RESET OPTIONALITY;
      };
      CREATE PROPERTY verbatim: std::str;
      CREATE REQUIRED PROPERTY year: std::int32 {
          SET REQUIRED USING (<std::int32>{});
          CREATE CONSTRAINT std::min_value(1000);
      };
  };
  CREATE TYPE traits::QualitativeTrait {
      CREATE REQUIRED MULTI PROPERTY modalities: std::str;
      CREATE REQUIRED PROPERTY name: std::str;
  };
  CREATE TYPE traits::QualitativeTraitTag {
      CREATE REQUIRED LINK trait: traits::QualitativeTrait;
      CREATE REQUIRED PROPERTY tag: std::str;
  };
  ALTER SCALAR TYPE occurrence::QuantityType EXTENDING enum<Unknown, One, Several, Ten, Tens, Hundred>;
  DROP MODULE reference;
};
