CREATE MIGRATION m1hck7qhij72j4hkzvlvn54ydzvvrsmbnl53mxma3ndgnimxqnjfba
    ONTO m1qrh76duzrnlyzkqagufayq6ggj7w5cxhzalg6gefjouabjjwmokq
{
  CREATE EXTENSION graphql VERSION '1.0';
  ALTER TYPE people::Person {
      ALTER PROPERTY first_name {
          CREATE CONSTRAINT std::max_len_value(32);
          CREATE CONSTRAINT std::min_len_value(2);
      };
      ALTER PROPERTY last_name {
          CREATE CONSTRAINT std::max_len_value(32);
          CREATE CONSTRAINT std::min_len_value(2);
      };
  };
};
