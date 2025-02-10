CREATE MIGRATION m1qm23x4gbpwnqjw72h3htgve37llf7utfxqxllb44dimjzzulk2tq
    ONTO m1cxemgp3fpboexirbglii5gwft3nq2pstecikw7ce2j4ga2jq4vsa
{
  CREATE FUNCTION location::coords_from_json(data: std::json) ->  tuple<precision: location::CoordinatesPrecision, latitude: std::float32, longitude: std::float32> USING (std::assert_exists((
      precision := <location::CoordinatesPrecision>(data)['precision'],
      latitude := <std::float32>(data)['latitude'],
      longitude := <std::float32>(data)['longitude']
  )));
  ALTER FUNCTION location::insert_site(data: std::json) USING (INSERT
      location::Site
      {
          name := <std::str>(data)['name'],
          code := <std::str>(data)['code'],
          description := <std::str>std::json_get(data, 'description'),
          coordinates := location::coords_from_json((data)['coordinates']),
          locality := <std::str>std::json_get(data, 'locality'),
          country := location::find_country(<std::str>std::json_get(data, 'country')),
          altitude := <std::int32>std::json_get(data, 'altitude')
      });
};
