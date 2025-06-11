CREATE MIGRATION m1dixofchznqbj7vqujkjwpszrcanaumyxa7l4slwthbclfpn4sh2a
    ONTO m1dikw6komtq5cnna6zqflk54353gkmxkseo5iiix3xoj4dg47uzyq
{
  ALTER TYPE location::Site {
      ALTER PROPERTY last_visited {
          USING (std::assert_single((WITH
              event_dates := 
                  .events.performed_on
          SELECT
              .events.performed_on FILTER
                  (.date = std::max(event_dates.date))
          LIMIT
              1
          )));
      };
  };
};
