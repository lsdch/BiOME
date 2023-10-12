CREATE MIGRATION m1cq5kdmkhjv26n2ndnguc36du5vv4kyvqq35jok7qngla2vdmv45a
    ONTO m1l5of6qlns35qygewok2pj4skrcx3ehty6f7p475r6abt3upvchpq
{
  ALTER TYPE location::Site {
      DROP LINK access_point;
      DROP LINK habitat;
  };
  DROP TYPE location::AccessPoint;
  DROP TYPE location::Habitat;
  CREATE TYPE location::HabitatTag EXTENDING default::Auditable {
      CREATE LINK parent: location::HabitatTag;
      CREATE PROPERTY color: std::str {
          CREATE REWRITE
              INSERT 
              USING (((.color ?? .parent.color) IF EXISTS (.parent) ELSE <std::str>{}));
          CREATE REWRITE
              UPDATE 
              USING (((.color ?? .parent.color) IF EXISTS (.parent) ELSE <std::str>{}));
      };
      CREATE PROPERTY description: std::str;
      CREATE PROPERTY is_required: std::bool {
          SET default := false;
      };
      CREATE REQUIRED PROPERTY label: std::str;
  };
  ALTER TYPE location::Site {
      CREATE REQUIRED MULTI LINK habitat_tags: location::HabitatTag {
          ON TARGET DELETE RESTRICT;
          SET REQUIRED USING (<location::HabitatTag>{});
          CREATE ANNOTATION std::title := 'A list of descriptors for the habitat that was targeted.';
      };
  };
};
