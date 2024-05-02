CREATE MIGRATION m1yy6mfdzjzrbr4iosyhxz5akkoce4x2sov3lmwzabgfvp2is4xirq
    ONTO m1xp3ir4lrxok6fgvk65sglcphxaigsgjp2sngxgand7mmkehw5ola
{
  ALTER TYPE location::HabitatTag {
      DROP PROPERTY color;
  };
  ALTER TYPE location::HabitatTag {
      DROP LINK parent;
  };
  ALTER TYPE location::HabitatTag {
      DROP PROPERTY is_required;
  };
  ALTER TYPE location::HabitatTag RENAME TO location::Habitat;
  ALTER TYPE location::Habitat {
      CREATE MULTI LINK depends: location::Habitat;
      CREATE MULTI LINK incompatibleFrom: location::Habitat;
      CREATE MULTI LINK incompatible := ((.incompatibleFrom UNION .<incompatibleFrom[IS location::Habitat]));
      ALTER PROPERTY label {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  ALTER TYPE location::Site {
      ALTER LINK habitat_tags {
          ON TARGET DELETE ALLOW;
      };
  };
};
