CREATE MIGRATION m12xvo3lo4db2q6c7tkxvwu5hhpfkhidiwwwzfy3cnx4r5kiynwveq
    ONTO m13vvofz3vmvdrugzv6woqkh3fl7sj6byrm7v5gsii3qaq4xx66dta
{
  ALTER TYPE events::AbioticMeasurement {
      ALTER LINK param {
          ON TARGET DELETE DELETE SOURCE;
      };
  };
};
