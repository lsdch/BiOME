CREATE MIGRATION m1sfktqdvf2xqrwx6af2sb5fs53vde2jcmih5t4ah32ica6epr5ztq
    ONTO m1qpq5lt4wt2e55vwvmxrtj36enpejzhywvigg5ujpftcemmxdcxxa
{
  ALTER TYPE occurrence::Identification {
      ALTER PROPERTY addendum {
          ALTER ANNOTATION std::description := "Additional information about the identification, e.g. 'group venustus'.";
      };
      CREATE PROPERTY confer: std::bool {
          SET default := false;
      };
  };
};
