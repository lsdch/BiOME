CREATE MIGRATION m15czjri44kvqhifef3m5p53s3roabgaozwtpocw7hloxswcmnia6q
    ONTO m14xzezsvdeffja6ewwjv4ngugttyli6oifdiihuqb7dzv3cwueoba
{
  ALTER TYPE occurrence::Occurrence {
      DROP LINK collections;
  };
  ALTER TYPE occurrence::Occurrence {
      CREATE MULTI PROPERTY collections: tuple<name: std::str, vouchers: array<std::str>> {
          CREATE ANNOTATION std::description := 'Collections where the specimen(s) can be found, with optional voucher identifiers.';
      };
  };
  DROP TYPE references::Collection;
};
