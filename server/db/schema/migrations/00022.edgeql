CREATE MIGRATION m1ez5o5finxnvighgosksewvxafspvlquayl4zj776ve45edviklpq
    ONTO m1pjt3ak4avjf6ktcjqsxctuk342bs2p7zz6hmnjgtaxm6gc7wslwa
{
  ALTER TYPE default::Auditable {
      ALTER LINK meta {
          DROP REWRITE
              UPDATE ;
          };
      };
  ALTER TYPE default::Auditable {
      ALTER LINK meta {
          CREATE REWRITE
              UPDATE 
              USING (SELECT
                  (UPDATE
                      .meta
                  SET {
                      modified := std::datetime_of_statement()
                  }) {
                      *
                  }
              );
      };
  };
};
