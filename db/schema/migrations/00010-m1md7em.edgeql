CREATE MIGRATION m1md7emgbcxwcqirm5pmsrdyndreyz2q5mrt7wrkw47e6awahpikya
    ONTO m1gjpyjnz7msaosgggqrewhrgbt2z4ey6gpvgbj4ske2q2qxn5c3zq
{
  DROP FUNCTION location::insert_site(data: std::json, infer_country: std::bool);
  CREATE FUNCTION location::insert_site(data: std::json, infer_country: OPTIONAL std::bool = false) ->  location::Site USING (WITH
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
          country := (IF infer_country THEN location::site_coords_to_country(coords) ELSE location::find_country(<std::str>std::json_get(data, 'country_code'))),
          altitude := <std::int32>std::json_get(data, 'altitude')
      });
};
