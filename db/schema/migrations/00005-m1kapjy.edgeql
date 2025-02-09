CREATE MIGRATION m1kapjymgvng52lp4fx4nvmq32kstmerhpcxvvpao2m5j6ouhben6q
    ONTO m1b5jq7vilsukd5c6wzl5dyhqn2mm3qtj3gynrlg7v2mi7vcpcltva
{
  CREATE FUNCTION people::insert_institution(data: std::json) ->  people::Institution USING (INSERT
      people::Institution
      {
          name := <std::str>(data)['name'],
          code := <std::str>(data)['code'],
          kind := <people::InstitutionKind>(data)['kind'],
          description := <std::str>std::json_get(data, 'description')
      });
};
