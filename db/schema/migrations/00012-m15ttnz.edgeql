CREATE MIGRATION m15ttnzxsxglzsk7lpbwjdzrk4pumpw76hrt7qm7etaom5ubyjlpeq
    ONTO m1c52quj6kiab4vp6xl5efdwlmrssvcxpzwq6mphvzg2zitcnj2gta
{
  ALTER TYPE events::Sampling {
      ALTER PROPERTY code {
          RESET EXPRESSION;
          RESET CARDINALITY;
          CREATE CONSTRAINT std::exclusive;
          SET TYPE std::str;
      };
      DROP PROPERTY suffix_id;
  };
  DROP FUNCTION events::sampling_code(s: events::Sampling, i: std::int32);
  CREATE FUNCTION events::sampling_code(e: events::Event, i: std::int32) ->  std::str USING (WITH
      number_suffix := 
          (<std::int32>(((SELECT
              events::SamplingCodeIndex
          FILTER
              (((.site = e.site) AND (.year = <std::int32>std::datetime_get(e.performed_on.date, 'year'))) AND (.month = <std::int32>std::datetime_get(e.performed_on.date, 'month')))
          )).count ?? 1) + i)
      ,
      date_suffix := 
          (SELECT
              std::assert_single((IF (e.performed_on.precision = date::DatePrecision.Unknown) THEN 'UNK' ELSE (IF (e.performed_on.precision = date::DatePrecision.Year) THEN <std::str>std::datetime_get(e.performed_on.date, 'year') ELSE (<std::str>std::datetime_get(e.performed_on.date, 'year') ++ std::str_pad_start(<std::str>std::datetime_get(e.performed_on.date, 'month'), 2, '0')))))
          )
      ,
      base := 
          std::assert_single(((e.site.code ++ '_') ++ date_suffix))
  SELECT
      std::assert_exists(std::assert_single(((base ++ '.') ++ <std::str>number_suffix)))
  );
};
