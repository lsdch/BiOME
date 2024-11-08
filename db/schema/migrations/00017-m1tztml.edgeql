CREATE MIGRATION m1tztmleuhdd3jz7si5vvhc5uvmi2xfbbg775zopfrspjg4zccgxea
    ONTO m1rxm2zur2jrfwsffhrnurhg3v73k5i2apetd3gbepe5k2nmcmxt5a
{
  ALTER TYPE events::Sampling {
      DROP PROPERTY is_donation;
  };
};
