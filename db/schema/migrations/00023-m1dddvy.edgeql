CREATE MIGRATION m1dddvy76mcxvqezphfwxkwzaz2rrri4ah7namlhwfwkee7ffmkoja
    ONTO m14dwoddp63uns2y62hjnc4xi4wey2spmfem2bquzfnytl67jc657a
{
  CREATE TYPE sampling::HabitatGroup EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY exclusive_elements: std::bool {
          SET default := true;
      };
      CREATE LINK depends: sampling::Habitat;
      CREATE REQUIRED PROPERTY label: std::str {
          CREATE CONSTRAINT std::exclusive;
      };
  };
  ALTER TYPE sampling::Habitat {
      CREATE REQUIRED LINK in_group: sampling::HabitatGroup {
          ON TARGET DELETE DELETE SOURCE;
          SET REQUIRED USING (<sampling::HabitatGroup>{});
      };
  };
  ALTER TYPE sampling::HabitatGroup {
      CREATE LINK elements := (.<in_group[IS sampling::Habitat]);
  };
  ALTER TYPE sampling::Habitat {
      CREATE LINK incompatible := ((.in_group.elements IF .in_group.exclusive_elements ELSE {}));
  };
};
