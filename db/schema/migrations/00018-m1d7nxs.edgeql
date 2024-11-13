CREATE MIGRATION m1d7nxsburgqsiaesq3guxncqorelaxa76q7js6zhda3wr4iynd4na
    ONTO m1tztmleuhdd3jz7si5vvhc5uvmi2xfbbg775zopfrspjg4zccgxea
{
  ALTER TYPE events::Event {
      DROP TRIGGER event_update;
  };
  ALTER TYPE events::Sampling {
      DROP TRIGGER index_insert;
  };
  DROP TYPE events::SamplingCodeIndex;
  ALTER TYPE reference::Article {
      ALTER PROPERTY authors {
          SET MULTI;
      };
      CREATE PROPERTY doi: std::str;
      CREATE PROPERTY journal: std::str;
  };
  ALTER TYPE reference::Article {
      ALTER PROPERTY year {
          DROP CONSTRAINT std::min_value(1500);
      };
  };
  ALTER TYPE reference::Article {
      ALTER PROPERTY year {
          CREATE CONSTRAINT std::min_value(1000);
      };
  };
};
