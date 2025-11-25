CREATE MIGRATION m1g5awc6ntbsknxzeuwwvzbqxohbyla7fkaw3sgeuvavybtike7c3a
    ONTO m1jyx4mftqlk3aqs5nlg7bt3o2b2qyxbd23g5bba7c33mhrzwrvida
{
  CREATE TYPE references::Collection EXTENDING default::Auditable, default::Vocabulary {
      CREATE PROPERTY contact: std::str;
  };
  ALTER TYPE occurrence::Occurrence {
      CREATE MULTI LINK collections: references::Collection {
          CREATE PROPERTY vouchers: array<std::str>;
      };
  };
  ALTER TYPE occurrence::Occurrence {
      DROP PROPERTY in_collection;
  };
  ALTER TYPE occurrence::Occurrence {
      DROP PROPERTY item_vouchers;
  };
};
