CREATE MIGRATION m1m4k5casiuunmop6wytotorrbapq3eicgk66jbzmjcuged3tblpfq
    ONTO m1bobxrlnufvwgoqzhfrej5ychutbfmrnvradd6h6zo2rovs4tt7jq
{
  ALTER TYPE occurrence::Occurrence {
      CREATE MULTI LINK datasets: datasets::OccurrenceDataset {
          ON SOURCE DELETE ALLOW;
          ON TARGET DELETE ALLOW;
      };
      ALTER LINK sampling {
          ON TARGET DELETE DELETE SOURCE;
      };
  };
  ALTER TYPE datasets::OccurrenceDataset {
      ALTER LINK occurrences {
          USING (.<datasets[IS occurrence::Occurrence]);
          RESET ON TARGET DELETE;
      };
  };
  ALTER TYPE seq::Sequence {
      ALTER LINK referenced_in {
          ON TARGET DELETE ALLOW;
      };
  };
  ALTER TYPE seq::ExternalSequence {
      ALTER LINK biomat {
          ON TARGET DELETE DELETE SOURCE;
      };
  };
};
