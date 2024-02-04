CREATE MIGRATION m13jebbrxp7su46hw6mqss7zkxpfvmkp7n7iez5j6slakjqs6sll2q
    ONTO m1pbbakrwhotrofbznchxbeb3m37t76tycc7odnyveeasafy3kkaqq
{
  ALTER TYPE default::Auditable {
      ALTER LINK meta {
          DROP REWRITE
              INSERT ;
          };
      };
  ALTER TYPE default::Auditable {
      ALTER LINK meta {
          CREATE REWRITE
              INSERT 
              USING (INSERT
                  default::Meta
              );
      };
  };
};
