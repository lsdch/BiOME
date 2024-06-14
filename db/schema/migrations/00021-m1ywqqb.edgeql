CREATE MIGRATION m1ywqqbqozldinqz7wci3rd3biuxpa725v354m565gfv27q5exvwua
    ONTO m1mjwfnafilprirzrjnnxjys7doq7wilphuwbji7fwia5w2dgl7k2q
{
  ALTER TYPE location::Country {
      DROP ANNOTATION std::description;
      DROP PROPERTY code;
      DROP PROPERTY name;
  };
  ALTER TYPE location::Locality {
      DROP LINK country;
      DROP CONSTRAINT std::exclusive ON ((.region, .municipality));
      DROP PROPERTY municipality;
      DROP PROPERTY region;
  };
  DROP TYPE location::Country;
  CREATE TYPE location::Country {
      CREATE ANNOTATION std::description := 'Countries as defined in the ISO 3166-1 norm.';
      CREATE REQUIRED PROPERTY code: std::str {
          CREATE CONSTRAINT std::exclusive;
          CREATE CONSTRAINT std::max_len_value(2);
          CREATE CONSTRAINT std::min_len_value(2);
      };
      CREATE REQUIRED PROPERTY name: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  ALTER TYPE location::Site {
      CREATE REQUIRED LINK country: location::Country {
          SET REQUIRED USING (<location::Country>{});
      };
  };
  ALTER TYPE location::Site {
      DROP LINK locality;
  };
  DROP TYPE location::Locality;
  ALTER TYPE location::Site {
      CREATE PROPERTY municipality: std::str;
  };
  ALTER TYPE location::Site {
      CREATE PROPERTY region: std::str;
  };
};
