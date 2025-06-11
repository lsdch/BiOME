CREATE MIGRATION m1gbrpj6kazku5fhety4tcujm7dxy7q7cnv2wvbnwqpuw5lehl65fq
    ONTO m1dixofchznqbj7vqujkjwpszrcanaumyxa7l4slwthbclfpn4sh2a
{
  ALTER TYPE location::Site {
      ALTER PROPERTY last_visited {
          USING (std::assert_single((SELECT
              .events.performed_on FILTER
                  (.date = std::max(__source__.events.performed_on.date))
          LIMIT
              1
          )));
      };
  };
};
