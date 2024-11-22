CREATE MIGRATION m1ok2r24epmn57wmjwloel7wqsr3dqtxcvmcvsyodnatqllruzn4wq
    ONTO m1dddvy76mcxvqezphfwxkwzaz2rrri4ah7namlhwfwkee7ffmkoja
{
  ALTER TYPE location::Site {
      ALTER LINK imported_in {
          DROP REWRITE
              INSERT ;
          };
      };
  ALTER TYPE location::SiteDataset RENAME TO datasets::Dataset;
  ALTER TYPE location::Site {
      ALTER LINK imported_in {
          CREATE REWRITE
              INSERT 
              USING ((IF (std::count(.datasets) = 1) THEN std::assert_single(.datasets) ELSE <datasets::Dataset>{}));
      };
  };
  ALTER TYPE occurrence::OccurrenceReport {
      CREATE PROPERTY in_collection: std::str;
  };
  ALTER TYPE occurrence::OccurrenceReport {
      CREATE MULTI PROPERTY item_voucher: std::str;
      DROP PROPERTY voucher;
  };
};
