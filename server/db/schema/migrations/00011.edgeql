CREATE MIGRATION m1pbbakrwhotrofbznchxbeb3m37t76tycc7odnyveeasafy3kkaqq
    ONTO m1p2557jlae6gqt772o7z7sbkaear6w3gx2hprvcgapo3faywfiqla
{
  ALTER TYPE default::Auditable {
      ALTER LINK meta {
          RESET default;
          CREATE REWRITE
              INSERT 
              USING (INSERT
                  default::Meta
                  {
                      created := std::datetime_of_statement()
                  });
      };
  };
};
