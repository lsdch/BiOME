CREATE MIGRATION m1vhuzbfmksam34as2anx5jomwdjlmzo5wf5uav4tzaqjnssefkgzq
    ONTO m1efr2jhr4mta6k7c6rz5hntfmlvn6p6evlh2bshftnevgnhiwxcpq
{
  CREATE SCALAR TYPE occurrence::IdentificationQualifier EXTENDING enum<CF, AFF>;
  ALTER TYPE occurrence::Identification {
      CREATE PROPERTY qualifier: occurrence::IdentificationQualifier {
          CREATE ANNOTATION std::description := "Identification qualifier, e.g. 'cf.' or 'aff.'";
      };
  };
};
