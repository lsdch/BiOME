CREATE MIGRATION m14dkru6mctwrtyf65fswdkahynnjtef2axylmyxoadhytuyzsihia
    ONTO m1yljweecpms65j2ngtg6cgk2tihbf5bm72nukz7uas5kgvjjr5tjq
{
  ALTER TYPE occurrence::Identification {
      ALTER PROPERTY identified_by {
          RESET CARDINALITY USING (SELECT
              .identified_by 
          LIMIT
              1
          );
      };
  };
};
