CREATE MIGRATION m13vvofz3vmvdrugzv6woqkh3fl7sj6byrm7v5gsii3qaq4xx66dta
    ONTO m1o4cxv25edjmtwx75delzovumzp7tmbfnwvpb6f3pje5bfly7o6ha
{
  ALTER FUNCTION date::from_json_with_precision(value: std::json) USING (std::assert_exists((
      date := (IF EXISTS ((value)['date']) THEN std::to_datetime(<std::int64>std::json_get(value, 'date', 'year'), <std::int64>std::json_get(value, 'date', 'month'), <std::int64>std::json_get(value, 'date', 'day'), 0, 0, 0, 'UTC') ELSE <std::datetime>{}),
      precision := <date::DatePrecision>(value)['precision']
  )));
};
