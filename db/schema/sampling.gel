module sampling {
  type Habitat extending default::Auditable {
    required label: str {
      constraint exclusive;
    };
    index on (.label);

    description: str;

    required in_group: HabitatGroup {
      on target delete delete source;
    };

    incompatible := (.in_group.elements if .in_group.exclusive_elements else {});
  }

  type HabitatGroup extending default::Auditable {
    required label: str {
      constraint exclusive;
    };
    index on (.label);

    depends: Habitat;
    required exclusive_elements: bool {
      default := true;
    };
    elements := .<in_group[is Habitat];
  }
}