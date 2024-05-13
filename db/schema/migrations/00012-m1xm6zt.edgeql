CREATE MIGRATION m1xm6ztndgvghtb7vh6i3fwxpfla3a3afgtmdgicimp423zp7wxmeq
    ONTO m1edi3ryzb25mbzbpxhh25fme4n45r5uo2djmtrsxn3mjwvyinguua
{
  ALTER TYPE location::Habitat {
      DROP LINK depends;
  };
  CREATE TYPE location::HabitatGroup EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY exclusive_elements: std::bool {
          SET default := true;
      };
      CREATE LINK depends: location::Habitat;
      CREATE REQUIRED PROPERTY label: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  ALTER TYPE location::Habitat {
      CREATE REQUIRED LINK in_group: location::HabitatGroup {
          SET REQUIRED USING (<location::HabitatGroup>{});
      };
  };
  ALTER TYPE location::HabitatGroup {
      CREATE LINK elements := (.<in_group[IS location::Habitat]);
  };
  ALTER TYPE location::Habitat {
      ALTER LINK incompatible {
          USING (((.incompatibleFrom UNION .<incompatibleFrom[IS location::Habitat]) UNION (.in_group.elements IF .in_group.exclusive_elements ELSE {})));
      };
  };
};
