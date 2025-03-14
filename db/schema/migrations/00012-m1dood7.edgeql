CREATE MIGRATION m1dood75z46qkztfemwnx6ibzvowqxvlj3y53rrnafgjyct2lrzmsa
    ONTO m1pxnh3jkudkehzohdcduh47ebzplqgng665u3zn3rhxtc6w5dkria
{
  ALTER TYPE location::Country {
      ALTER PROPERTY code {
          CREATE CONSTRAINT std::min_len_value(3);
      };
  };
  ALTER TYPE location::Country {
      ALTER PROPERTY code {
          DROP CONSTRAINT std::min_len_value(2);
      };
      CREATE REQUIRED PROPERTY continent: std::str {
          SET REQUIRED USING (<std::str>{});
      };
      CREATE REQUIRED PROPERTY subcontinent: std::str {
          SET REQUIRED USING (<std::str>{});
      };
  };
};
