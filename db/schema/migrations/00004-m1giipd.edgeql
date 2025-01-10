CREATE MIGRATION m1giipdbqfje4e4ktps6mah36vtlan6kq3hwk6ubxyrd6uselpt6ya
    ONTO m1sd2ildlby4zjxatkwk5mxpcdwhgkgcnoom63dpt37aiuu7wzswnq
{
  CREATE ALIAS occurrence::BioMaterialWithType := (
      SELECT
          occurrence::BioMaterial {
              required has_sequences := EXISTS (([IS occurrence::ExternalBioMat].sequences ?? [IS occurrence::InternalBioMat].specimens.sequences)),
              required is_homogenous := ([IS occurrence::ExternalBioMat].homogenous ?? ([IS occurrence::InternalBioMat].homogenous ?? true)),
              required is_congruent := ([IS occurrence::ExternalBioMat].congruent ?? ([IS occurrence::InternalBioMat].congruent ?? true)),
              sequence_consensus := ([IS occurrence::ExternalBioMat].seq_consensus ?? [IS occurrence::InternalBioMat].seq_consensus)
          }
  );
};
