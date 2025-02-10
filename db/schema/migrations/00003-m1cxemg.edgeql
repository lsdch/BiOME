CREATE MIGRATION m1cxemgp3fpboexirbglii5gwft3nq2pstecikw7ce2j4ga2jq4vsa
    ONTO m1z4w7raqyomlkgxtyi6vt25aoexfdhiyedzlqjhwzka7zrzml32pq
{
  CREATE FUNCTION location::find_country(code: OPTIONAL std::str) -> OPTIONAL location::Country USING ((IF EXISTS (code) THEN (SELECT
      std::assert_exists(location::Country FILTER
          (.code = code)
      , message := ('Invalid country code: ' ++ code))
  ) ELSE <location::Country>{}));
  CREATE FUNCTION location::insert_site(data: std::json) ->  location::Site USING (INSERT
      location::Site
      {
          name := <std::str>(data)['name'],
          code := <std::str>(data)['code'],
          description := <std::str>std::json_get(data, 'description'),
          coordinates := (
              precision := <location::CoordinatesPrecision>(data)['coordinates']['precision'],
              latitude := <std::float32>(data)['coordinates']['latitude'],
              longitude := <std::float32>(data)['coordinates']['longitude']
          ),
          locality := <std::str>std::json_get(data, 'locality'),
          country := location::find_country(<std::str>std::json_get(data, 'country')),
          altitude := <std::int32>std::json_get(data, 'altitude')
      });
};
