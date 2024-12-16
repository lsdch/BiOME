CREATE MIGRATION m1rllj7uzgqd2dco2gymr4e7k63za7atfp3ioygpttnxh27ijrjyea
    ONTO m12at5dpt2wnfmrg4ngtk3fwbidnzcad6iblfh5bvydafaz5cml3aa
{
  ALTER TYPE occurrence::BioMaterial {
      CREATE REQUIRED PROPERTY category := (std::assert_exists((IF (__source__ IS occurrence::InternalBioMat) THEN occurrence::OccurrenceCategory.Internal ELSE (IF (__source__ IS occurrence::ExternalBioMat) THEN occurrence::OccurrenceCategory.External ELSE <occurrence::OccurrenceCategory>{})), message := (('Occurrence category for occurrence::BioMaterial subtype ' ++ __source__.__type__.name) ++ ' is undefined')));
  };
};
