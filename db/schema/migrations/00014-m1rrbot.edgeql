CREATE MIGRATION m1rrbotd7anfquiocfntsyz5mbxz7gs6i7hukr5lcbmappxkng65nq
    ONTO m1fovhcwjrs6rznqorfxdso5weft3nr3ntnbqrrjviahz2vpybnziq
{
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
              (IF EXISTS (std::json_get(data, 'country_code')) THEN location::find_country(<std::str>std::json_get(data, 'country_code')) ELSE std::assert_exists(location::site_coords_to_country(coords), message := ('Country not found for the given coordinates' ++ std::to_str(<std::json>coords))))
          ),
          altitude := <std::int32>std::json_get(data, 'altitude')
      });
};
