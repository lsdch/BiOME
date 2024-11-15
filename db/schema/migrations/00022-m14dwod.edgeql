CREATE MIGRATION m14dwoddp63uns2y62hjnc4xi4wey2spmfem2bquzfnytl67jc657a
    ONTO m1ktchkrmzvi33nnju4pes4laahsrlxpypklkacd5zzfquqb3vbrwa
{
  CREATE MODULE sampling IF NOT EXISTS;
  ALTER TYPE location::HabitatGroup {
      DROP LINK depends;
  };
  ALTER TYPE location::Habitat {
      DROP LINK incompatible;
  };
  ALTER TYPE location::HabitatGroup {
      DROP LINK elements;
      DROP PROPERTY exclusive_elements;
      DROP PROPERTY label;
  };
  ALTER TYPE location::Habitat {
      DROP LINK in_group;
  };
  ALTER TYPE location::Habitat {
      DROP LINK incompatible_from;
  };
  ALTER TYPE location::Habitat RENAME TO sampling::Habitat;
  DROP TYPE location::HabitatGroup;
};
