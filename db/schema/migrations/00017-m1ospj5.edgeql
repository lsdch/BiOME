CREATE MIGRATION m1ospj5bhcqtnqm3pvzrhoummiq5jyk7dhfibxz46njagvhkqy3cba
    ONTO m17l2kvzammg3y6zxoosit6b2sb3pjxsg2achldffevzjlhnop6vzq
{
  ALTER FUNCTION location::site_coords_to_country(coords: tuple<precision: location::CoordinatesPrecision, latitude: std::float32, longitude: std::float32>) USING (SELECT
      location::position_to_country(coords.latitude, coords.longitude)
  );
  ALTER FUNCTION location::insert_site(data: std::json) USING (WITH
      coords := 
          location::coords_from_json((data)['coordinates'])
  INSERT
      location::Site
      {
          name := <std::str>(data)['name'],
          code := <std::str>(data)['code'],
          description := <std::str>std::json_get(data, 'description'),
          coordinates := coords,
          locality := <std::str>std::json_get(data, 'locality'),
          country := (SELECT
              (IF EXISTS (std::json_get(data, 'country_code')) THEN location::find_country(<std::str>(data)['country_code']) ELSE std::assert_exists(location::site_coords_to_country(coords), message := ('Country not found for the given coordinates: ' ++ std::to_str(data))))
          ),
          altitude := <std::int32>std::json_get(data, 'altitude')
      });
};
