CREATE MIGRATION m1efr2jhr4mta6k7c6rz5hntfmlvn6p6evlh2bshftnevgnhiwxcpq
    ONTO m14o4u2w6uykl5oido6eamuoamtgtmhs74k56nj77qf57ge7tgtiiq
{
  ALTER TYPE occurrence::ExternalBioMat {
      CREATE SINGLE LINK best_identification := (SELECT
          (.seq_consensus ?? .identification.taxon)
      );
  };
  ALTER TYPE events::Sampling {
      ALTER LINK occurring_taxa {
          USING (SELECT
              DISTINCT (((.samples[IS occurrence::ExternalBioMat].best_identification UNION ((SELECT
                  .external_seqs
              FILTER
                  NOT (EXISTS (.source_sample))
              )).identification.taxon) UNION .samples[IS occurrence::InternalBioMat].identified_taxa))
          );
      };
  };
};
