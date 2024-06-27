CREATE MIGRATION m1qwlpsc5dcohpkmplks5s7xbhrlxoklqlyonpzxd44wnp3sagkd5a
    ONTO m1ywqqbqozldinqz7wci3rd3biuxpa725v354m565gfv27q5exvwua
{
  ALTER TYPE location::Country {
      CREATE MULTI LINK sites := (.<country[IS location::Site]);
  };
  CREATE ALIAS location::CountryList := (
      SELECT
          location::Country {
              *,
              sites_count := std::count(.sites)
          }
  );
  ALTER TYPE location::Site {
      ALTER PROPERTY altitude {
          SET TYPE std::int32 USING (<std::int32>.altitude.min);
      };
  };
  ALTER TYPE location::Site {
      ALTER PROPERTY code {
          CREATE CONSTRAINT std::max_len_value(10);
      };
  };
  ALTER TYPE location::Site {
      ALTER PROPERTY code {
          DROP CONSTRAINT std::max_len_value(8);
      };
  };
  CREATE SCALAR TYPE location::CoordinatesPrecision EXTENDING enum<`<100m`, `<1Km`, `<10Km`, `10-100Km`, Unknown>;
  ALTER TYPE location::Site {
      ALTER PROPERTY coordinates {
          SET TYPE tuple<precision: location::CoordinatesPrecision, latitude: std::float32, longitude: std::float32> USING (<tuple<precision: location::CoordinatesPrecision, latitude: std::float32, longitude: std::float32>>{});
      };
  };
  ALTER TYPE location::Site {
      ALTER PROPERTY municipality {
          RENAME TO locality;
      };
  };
  ALTER TYPE location::Site {
      CREATE PROPERTY access_point: std::str;
      DROP PROPERTY region;
  };
  DROP SCALAR TYPE location::CoordinateMaxPrecision;
};
