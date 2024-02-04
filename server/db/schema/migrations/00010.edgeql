CREATE MIGRATION m1p2557jlae6gqt772o7z7sbkaear6w3gx2hprvcgapo3faywfiqla
    ONTO m1s5gn5r34j6ya2dfz7wphacwt6multa5bkckjmujsyaa2tmaxcvbq
{
  ALTER TYPE default::Auditable {
      ALTER LINK meta {
          SET default := (INSERT
              default::Meta
              {
                  created := std::datetime_of_statement()
              });
      };
  };
};
