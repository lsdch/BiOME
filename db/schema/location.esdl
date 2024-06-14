
module location {

  type Country {
    annotation description := "Countries as defined in the ISO 3166-1 norm.";
    required name: str {
      constraint exclusive;
    };
    required code: str {
      constraint exclusive;
      constraint min_len_value(2);
      constraint max_len_value(2);
    };

    multi sites := .<country[is Site];
  }

  alias CountryList := (
    select Country { *, sites_count := count(.sites) }
  );

  type Habitat extending default::Auditable {
    required label: str {
      constraint exclusive;
    };
    description: str;

    required in_group: HabitatGroup {
      on target delete delete source;
    };

    multi incompatible_from: Habitat {
      on target delete allow;
    };
    incompatible := (
      .incompatible_from
      union .<incompatible_from[is Habitat]
      union (.in_group.elements if .in_group.exclusive_elements else {})
    )
  }

  type HabitatGroup extending default::Auditable {
    required label: str {
      constraint exclusive;
    };
    depends: Habitat;
    required exclusive_elements: bool {
      default := true;
    };
    elements := .<in_group[is Habitat];
  }

  scalar type CoordinateMaxPrecision extending enum<"m10", "m100", "Km1", "Km10", "Km100", "Unknown">;

  type Site extending default::Auditable {
    required name : str { constraint exclusive };
    required code : str {
      annotation title := "Site identifier";
      annotation description := "A short, unique, user-generated, alphanumeric identifier. Recommended size is 8.";
      constraint exclusive;
      constraint min_len_value(4);
      constraint max_len_value(8);
    }
    description: str;

    # required multi habitat_tags: Habitat {
    #   annotation title := "A list of descriptors for the habitat that was targeted.";
    #   on target delete allow;
    # };

    region: str;
    municipality: str;
    required country: Country;

    coordinates: tuple<
      precision: CoordinateMaxPrecision,
      latitude: float32,
      longitude: float32
    > {
      constraint expression on (
        (.latitude <= 90 and .latitude >= -90
        and .longitude <= 180 and .longitude >= -180)
      );
    };

    altitude: tuple<min: float32, max: float32> {
      annotation title := "The site elevation in meters";
    };


    multi link abiotic_measurements := .<site[is event::AbioticMeasurement];
    multi link samplings := .<site[is event::Sampling];
    multi link spottings := .<site[is event::Spotting];
    multi datasets := .<sites[is SiteDataset];
  }

  type SiteDataset extending default::Auditable {
    required label: str {
      constraint min_len_value(4);
      constraint max_len_value(40);
    }
    description: str;
    multi sites: Site;
    required multi maintainers: people::Person;
  }
}