CREATE MIGRATION m1dentl2ncfb43uknlst3ms6jsscjxo55iwjyknynb6js3ykblnvla
    ONTO m12pk4hd3qqzbuejeit7y6wkgru5qrqegt4fpvravcxrfgrpfxcosa
{
  ALTER TYPE default::AppConfig {
      ALTER PROPERTY public {
          CREATE ANNOTATION std::description := 'Sets whether parts of the platform are open to the public (anonymous users).';
      };
  };
  ALTER TYPE people::Institution {
      CREATE REQUIRED PROPERTY acronym: std::str {
          SET REQUIRED USING (<std::str>{});
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(12);
          CREATE CONSTRAINT std::min_len_value(2);
      };
  };
  ALTER TYPE people::Institution {
      ALTER PROPERTY comments {
          RENAME TO description;
      };
      ALTER PROPERTY name {
          CREATE CONSTRAINT std::max_len_value(128);
          CREATE CONSTRAINT std::min_len_value(10);
      };
  };
};
