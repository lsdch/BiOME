CREATE MIGRATION m17l2kvzammg3y6zxoosit6b2sb3pjxsg2achldffevzjlhnop6vzq
    ONTO m1kbfwewg4h67qwqh42y4wioz6ssirwqqmopu73ruq3qibwuuyy7ma
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
              (IF EXISTS (std::json_get(data, 'country_code')) THEN location::find_country(<std::str>(data)['country_code']) ELSE location::site_coords_to_country(coords))
          ),
          altitude := <std::int32>std::json_get(data, 'altitude')
      });
};
