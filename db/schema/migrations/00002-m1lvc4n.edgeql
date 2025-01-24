CREATE MIGRATION m1lvc4nif6vlxn23m7rtth7dpexhmsbqmk6g2gcsixuqobh7qstpta
    ONTO m1kxwncvxgkrc7odf4nmle27e46qiuhdfnta4zlfawvxzr3gedv4qa
{
  ALTER TYPE seq::AssembledSequence DROP EXTENDING default::CodeIdentifier;
  ALTER TYPE seq::ExternalSequence DROP EXTENDING default::CodeIdentifier;
  ALTER TYPE seq::ExternalSequence {
      ALTER PROPERTY code {
          RESET readonly;
          RESET CARDINALITY;
      };
  };
};
