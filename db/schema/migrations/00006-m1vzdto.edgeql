CREATE MIGRATION m1vzdto5xraqnpbal2qgyu3koxct5fwlrrw27g2xasahvemugvujhq
    ONTO m1kapjymgvng52lp4fx4nvmq32kstmerhpcxvvpao2m5j6ouhben6q
{
  CREATE FUNCTION people::insert_or_create_institution(data: std::json) ->  people::Institution USING ((IF (std::json_typeof(data) = 'object') THEN (SELECT
      people::insert_institution(data)
  ) ELSE (IF (std::json_typeof(data) = 'string') THEN (SELECT
      std::assert_exists(people::Institution FILTER
          (.code = <std::str>data)
      , message := ('Failed to find institution with code: ' ++ <std::str>data))
  ) ELSE std::assert_exists(<people::Institution>{}, message := ('Invalid institution JSON type: ' ++ std::json_typeof(data))))));
  CREATE FUNCTION people::insert_person(data: std::json) ->  people::Person USING (INSERT
      people::Person
      {
          first_name := <std::str>(data)['first_name'],
          last_name := <std::str>(data)['last_name'],
          alias := <std::str>std::json_get(data, 'alias'),
          contact := <std::str>std::json_get(data, 'contact'),
          comment := <std::str>std::json_get(data, 'comment'),
          institutions := DISTINCT ((FOR inst IN std::json_array_unpack(std::json_get(data, 'institutions'))
          UNION 
              people::insert_or_create_institution(inst)))
      });
};
