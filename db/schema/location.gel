
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

  function find_country(code: optional str) -> optional Country {
    using (
      if exists code then (
        select assert_exists(
          location::Country
          filter .code = code,
          message := ("Invalid country code: " ++ code)
        )
      ) else <location::Country>{}
    )
  }

  alias CountryList := (
    select Country { *, sites_count := count(.sites) }
  );


  # <100m: Coordinates of site position
  # <1Km: Nearest small locality
  # <10Km: Nearest locality
  # 10-100Km: Country/region
  # Unknown
  scalar type CoordinatesPrecision extending enum<"<100m", "<1km", "<10km", "10-100km", "Unknown">;

  type Site extending default::Auditable {
    required name : str;
    required code : str {
      annotation title := "Site identifier";
      annotation description := "A short, unique, user-generated, alphanumeric identifier. Recommended size is 8.";
      constraint exclusive;
      constraint min_len_value(3);
      constraint max_len_value(10);
    };

    index on (.code);

    description: str;

    locality: str;
    country: Country;
    required user_defined_locality: bool {
      default := false
    };

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
    multi datasets := .<sites[is datasets::SiteDataset];

    # imported_in: datasets::AbstractDataset {
    #   on target delete allow;
    #   on source delete allow;
    # };
  }

  function coords_from_json(data: json) -> tuple<precision: CoordinatesPrecision, latitude: float32, longitude: float32> {
    using (
      assert_exists((
        precision := <CoordinatesPrecision>data['precision'],
        latitude := <float32>data['latitude'],
        longitude := <float32>data['longitude']
      ))
    )
  }

  function insert_site(data: json) -> Site {
    using (insert Site {
      name := <str>data['name'],
      code := <str>data['code'],
      description := <str>json_get(data, 'description'),
      coordinates := coords_from_json(data['coordinates']),
      locality := <str>json_get(data, 'locality'),
      country := find_country(<str>json_get(data, 'country_code')),
      altitude := <int32>json_get(data, 'altitude'),
    })
  }
}