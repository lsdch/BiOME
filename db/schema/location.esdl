
module location {

  type Country extending default::Auditable {
    annotation description := "Countries as defined in the ISO 3166-1 norm.";
    required name: str {
      constraint exclusive;
      readonly := true;
    };
    required code: str {
      constraint exclusive;
      readonly := true;
      constraint min_len_value(2);
      constraint max_len_value(2);
    };
    multi link localities := .<country[is Locality]
  }

  type Locality extending default::Auditable {
    region: str;
    municipality: str;
    constraint exclusive on ((.region, .municipality));

    required code: str {
      constraint exclusive;
      annotation description := "Format like 'municipality|region[country_code]'";
      rewrite insert, update using (
        .municipality ?? "NA" ++ "|"
        ++ .region ?? "NA" ++ "|"
        ++ .country.code
      )
    };

    required country: Country;
    multi link sites := .<locality[is Site]
  }


  type Habitat extending default::Auditable {
    required label: str {
      constraint exclusive;
    };
    description: str;
    multi depends: Habitat;
    multi incompatibleFrom: Habitat;
    incompatible := (.incompatibleFrom union .<incompatibleFrom[is Habitat])
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

    required multi habitat_tags: Habitat {
      annotation title := "A list of descriptors for the habitat that was targeted.";
      on target delete allow;
    };

    required locality: Locality {
      on target delete restrict;
    };

    required coordinates: tuple<
      precision: CoordinateMaxPrecision,
      latitude: float32,
      longitude: float32
    > {
      constraint expression on (
        (.latitude <= 90 and .latitude >= -90
        and .longitude <= 180 and .longitude >= -180)
      );
    };

    altitudeRange: tuple<min:float32, max:float32> {
      annotation title := "The site elevation in meters";
    };


    multi link abiotic_measurements := .<site[is event::AbioticMeasurement];
    multi link samplings := .<site[is event::Sampling];
    multi link spottings := .<site[is event::Spotting];
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