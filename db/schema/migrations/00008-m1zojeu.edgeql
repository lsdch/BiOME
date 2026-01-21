CREATE MIGRATION m1zojeutj47b5airchz4mwkcrzggkwmb2q3oat6vesmuhvyt4f27za
    ONTO m1bm6nwr4rz4s2byr7ybmzqljuc4vtfs6r4yom3kaarprpt3vrwnua
{
  ALTER TYPE location::Site {
      ALTER PROPERTY code {
          DROP CONSTRAINT std::max_len_value(10);
      };
  };
  ALTER TYPE location::Site {
      ALTER PROPERTY code {
          CREATE CONSTRAINT std::max_len_value(32);
      };
  };
};
