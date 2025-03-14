CREATE MIGRATION m1aabgqjxm74kmco3jgssfyalpgzi6mf4kqs2fogkccmqxqipiosrq
    ONTO m12lk2vyqr7ey6ddbd4tbzpqthemqjtc6kcmu53ojumrwsdkxqplpq
{
  ALTER TYPE location::Country {
      DROP PROPERTY boundaries;
  };
  CREATE TYPE location::CountryBoundary {
      CREATE REQUIRED LINK country: location::Country {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY geometry: ext::postgis::geometry {
          CREATE ANNOTATION std::description := 'PostGIS polygon/multipolygon defining the country boundary.';
      };
  };
};
