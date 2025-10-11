CREATE MIGRATION m1s526bxy4nvyiyin27obsd27ydaiifbv23dypniirvtr4xfq7exma
    ONTO m1f7bean47xanlluwcpovqiry3k3x6a6ajdoxjpfefk3maj7n3jcdq
{
  ALTER TYPE occurrence::BioMaterial {
      CREATE ACCESS POLICY no_duplicated_taxon_in_same_sampling
          DENY UPDATE, INSERT USING (EXISTS ((std::count(((GROUP
              .sampling.occurrences
          USING
              taxon := 
                  .identification.taxon
          BY taxon)).elements) > 1)));
  };
};
