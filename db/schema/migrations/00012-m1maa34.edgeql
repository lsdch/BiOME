CREATE MIGRATION m1maa34zvuro4hoybz3pfnev3mgj7pbp7tgoa7oicqmez3aod6jxka
    ONTO m1xn7awaoifsolgq55fpssj4556k4gkcd5476qhxxxysr65sy6fj7a
{
  ALTER TYPE location::Site {
      CREATE PROPERTY last_visited := ((std::assert_single((WITH
          event_dates := 
              .events.performed_on
      SELECT
          .events.performed_on FILTER
              (.date = std::max(event_dates.date))
      LIMIT
          1
      )),));
  };
};
