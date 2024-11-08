CREATE MIGRATION m1rxm2zur2jrfwsffhrnurhg3v73k5i2apetd3gbepe5k2nmcmxt5a
    ONTO m13nmsmvfnx6bnalbpk245ea2fupkbamdua7csqyfwzbz7ij73xa4a
{
  ALTER TYPE events::Sampling {
      CREATE MULTI LINK habitats: location::Habitat;
      CREATE MULTI PROPERTY access_points: std::str;
  };
  ALTER TYPE location::Site {
      DROP PROPERTY access_point;
  };
};
