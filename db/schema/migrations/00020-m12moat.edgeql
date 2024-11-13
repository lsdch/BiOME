CREATE MIGRATION m12moatuy6u3hxui2xrsqhctqep2uywc7kfpwu2iaoasjd2repfpha
    ONTO m1cyfkatidiapmm7kxfcpye4nvcpn6xf36exrafcecuv7osw7lsufa
{
  ALTER TYPE seq::Gene {
      DROP PROPERTY code;
  };
  ALTER TYPE seq::Gene {
      CREATE PROPERTY code: std::str {
          SET REQUIRED USING (<std::str>{});
      };
  };
  ALTER TYPE seq::Gene {
      CREATE PROPERTY label: std::str {
          SET REQUIRED USING (<std::str>{});
      };
      DROP PROPERTY name;
      EXTENDING default::Vocabulary BEFORE default::Auditable;
      ALTER PROPERTY code {
          RESET OPTIONALITY;
          DROP OWNED;
          RESET TYPE;
      };
      ALTER PROPERTY description {
          DROP OWNED;
          RESET TYPE;
      };
      ALTER PROPERTY motu {
          SET default := false;
      };
  };
  ALTER TYPE seq::Gene {
      ALTER PROPERTY label {
          RESET OPTIONALITY;
          DROP OWNED;
          RESET TYPE;
      };
  };
};
