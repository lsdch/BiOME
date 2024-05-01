CREATE MIGRATION m1xp3ir4lrxok6fgvk65sglcphxaigsgjp2sngxgand7mmkehw5ola
    ONTO m136pufzhmard4gteawe6idtkdaxhhtnzzxnku2sbdbabt2guscrkq
{
  ALTER TYPE admin::Settings {
      ALTER LINK email {
          USING (SELECT
              admin::EmailSettings 
          LIMIT
              1
          );
          RESET ON SOURCE DELETE;
      };
      ALTER LINK instance {
          USING (std::assert_exists((SELECT
              admin::InstanceSettings 
          LIMIT
              1
          )));
          RESET ON SOURCE DELETE;
          DROP CONSTRAINT std::exclusive;
      };
      ALTER LINK security {
          USING (std::assert_exists((SELECT
              admin::SecuritySettings 
          LIMIT
              1
          )));
          RESET ON SOURCE DELETE;
          DROP CONSTRAINT std::exclusive;
      };
  };
};
