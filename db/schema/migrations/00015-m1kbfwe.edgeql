CREATE MIGRATION m1kbfwewg4h67qwqh42y4wioz6ssirwqqmopu73ruq3qibwuuyy7ma
    ONTO m1rrbotd7anfquiocfntsyz5mbxz7gs6i7hukr5lcbmappxkng65nq
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
              (IF EXISTS (std::json_get(data, 'country_code')) THEN location::find_country(<std::str>std::json_get(data, 'country_code')) ELSE location::site_coords_to_country(coords))
          ),
          altitude := <std::int32>std::json_get(data, 'altitude')
      });
};
