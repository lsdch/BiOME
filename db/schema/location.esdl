
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

  # <100m: Coordinates of site position
  # <1Km: Nearest small locality
  # <10Km: Nearest locality
  # 10-100Km: Country/region
  # Unknown
  scalar type CoordinatesPrecision extending enum<"<100m", "<1Km", "<10Km", "10-100Km", "Unknown">;

  type Site extending default::Auditable {
    required name : str { constraint exclusive };
    required code : str {
      annotation title := "Site identifier";
      annotation description := "A short, unique, user-generated, alphanumeric identifier. Recommended size is 8.";
      constraint exclusive;
      constraint min_len_value(4);
      constraint max_len_value(10);
    }
    description: str;

    # required multi habitat_tags: Habitat {
    #   annotation title := "A list of descriptors for the habitat that was targeted.";
    #   on target delete allow;
    # };

    locality: str;
    required country: Country;

    required coordinates: tuple<
      precision: CoordinatesPrecision,
      latitude: float32,
      longitude: float32
    > {
      constraint expression on (
        (.latitude <= 90 and .latitude >= -90
        and .longitude <= 180 and .longitude >= -180)
      );
      rewrite insert, update using ((
        precision:= __subject__.coordinates.precision,
        latitude := <float32>round(<decimal>__subject__.coordinates.latitude, 5),
        longitude := <float32>round(<decimal>__subject__.coordinates.longitude, 5)
      ));
    };

    altitude: int32 {
      annotation title := "The site elevation in meters";
    };

    multi link events := .<site[is events::Event];
    multi datasets : SiteDataset {
      on target delete allow;
      on source delete allow;
    };

    imported_in: SiteDataset {
      on target delete allow;
      on source delete allow;
      rewrite insert using (
        if count(.datasets) = 1 then assert_single(.datasets) else <SiteDataset>{}
      );
    };
  }

  type SiteDataset extending default::Auditable {
    required label: str {
      constraint min_len_value(4);
      constraint max_len_value(40);
    }
    required slug: str {
      constraint exclusive;
    };

    description: str;
    multi link sites := .<datasets[is Site];
    required multi maintainers: people::Person {
      # edgedb error: SchemaDefinitionError:
      # cannot specify a rewrite for link 'maintainers' of object type 'location::SiteDataset' because it is multi
      # Hint: this is a temporary implementation restriction

      # rewrite insert using (global default::current_user.identity union .maintainers);
      # rewrite update using (.maintainers union .meta.created_by_user.identity);
    };
  }
}