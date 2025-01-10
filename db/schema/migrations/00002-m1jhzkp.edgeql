CREATE MIGRATION m1jhzkpurb6uv4sdhl7qzvfvlizdcchdylrbwlawno6ij4prkhe7la
    ONTO m1zl7heywj5fllwx2kcua3mmj45gzih7wihvmljorznpqmaaiug4vq
{
  ALTER ALIAS occurrence::BioMaterialWithType USING (SELECT
      occurrence::BioMaterial {
          required has_sequences := EXISTS (([IS occurrence::ExternalBioMat].sequences ?? [IS occurrence::InternalBioMat].specimens.sequences)),
          required is_homogenous := ([IS occurrence::ExternalBioMat].is_homogenous ?? ([IS occurrence::InternalBioMat].is_homogenous ?? true)),
          required is_congruent := ([IS occurrence::ExternalBioMat].is_congruent ?? ([IS occurrence::InternalBioMat].is_congruent ?? true)),
          sequence_consensus := ([IS occurrence::ExternalBioMat].sequence_consensus ?? [IS occurrence::InternalBioMat].sequence_consensus)
      }
  );
  ALTER TYPE samples::Specimen {
      ALTER PROPERTY molecular_number {
          DROP REWRITE
              INSERT ;
          };
      };
  ALTER TYPE samples::Specimen {
      ALTER PROPERTY molecular_number {
          CREATE REWRITE
              INSERT 
              USING (<std::str>(SELECT
                  std::count(samples::Specimen FILTER
                      (samples::Specimen.biomat.sampling.event.site = __subject__.biomat.sampling.event.site)
                  )
              ));
      };
  };
};
