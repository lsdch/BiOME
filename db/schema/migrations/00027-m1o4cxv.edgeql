CREATE MIGRATION m1o4cxv25edjmtwx75delzovumzp7tmbfnwvpb6f3pje5bfly7o6ha
    ONTO m1y7jjaj4wamvhqa6x3mebswhokj4wyavmczutonedxhgrzzxdwswa
{
  CREATE FUNCTION date::from_json_with_precision(value: std::json) ->  tuple<date: std::datetime, precision: date::DatePrecision> USING (std::assert_exists((
      date := (IF EXISTS ((value)['date']) THEN std::to_datetime(<std::int64>(value)['date']['year'], <std::int64>(value)['date']['month'], <std::int64>(value)['date']['day'], 0, 0, 0, 'UTC') ELSE <std::datetime>{}),
      precision := <date::DatePrecision>(value)['precision']
  )));
};
