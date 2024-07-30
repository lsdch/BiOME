CREATE MIGRATION m1b2dodwjgwhy7vfu5lsvvpc7njb35xtvtuuwk7jgmdezjegr4niza
    ONTO m15nns5vqagvqv5viond2vefn2fs4gvpqj35unxrwl6kq7e3wxgg6q
{
  ALTER TYPE event::Action {
      ALTER LINK event {
          ON TARGET DELETE DELETE SOURCE;
      };
  };
  ALTER TYPE event::Event {
      CREATE MULTI LINK abiotic_measurements := (.<event[IS event::AbioticMeasurement]);
  };
  ALTER TYPE event::Sampling {
      DROP LINK measurements;
  };
  ALTER TYPE event::AbioticMeasurement {
      DROP LINK related_sampling;
  };
  ALTER TYPE event::Event {
      DROP LINK actions;
  };
  ALTER TYPE event::Event {
      CREATE MULTI LINK samplings := (.<event[IS event::Sampling]);
      ALTER LINK site {
          ON SOURCE DELETE ALLOW;
          ON TARGET DELETE DELETE SOURCE;
      };
  };
  ALTER TYPE event::Event {
      CREATE MULTI LINK spottings := (.<event[IS event::Spotting]);
  };
  ALTER TYPE location::Site {
      ALTER LINK imported_in {
          ON SOURCE DELETE ALLOW;
          ON TARGET DELETE ALLOW;
      };
  };
  ALTER TYPE location::SiteDataset {
      ALTER LINK sites {
          ON SOURCE DELETE ALLOW;
          ON TARGET DELETE ALLOW;
      };
  };
};
