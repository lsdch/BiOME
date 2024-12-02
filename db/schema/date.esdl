module date {
  scalar type DatePrecision extending enum<Year, Month, Day, Unknown>;

  abstract constraint required_unless_unknown(date: datetime, precision: DatePrecision) {
    errmessage := "Date value is required, except when precision is 'Unknown'";
    using (exists date if precision != DatePrecision.Unknown else not exists date)
  }

  function rewrite_date(value: tuple<date:datetime, precision:DatePrecision>) -> optional datetime
  using (
      datetime_truncate(value.date, 'years') if value.precision = DatePrecision.Year else
      datetime_truncate(value.date, 'months') if value.precision = DatePrecision.Month else
      datetime_truncate(value.date, 'days') if value.precision = DatePrecision.Day else
      <datetime>{}
  );

  function from_json_with_precision(value: json) -> tuple<date:datetime, precision:DatePrecision>
  using (assert_exists((
    date := (
      if exists value['date'] then to_datetime(
        <int64>json_get(value, 'date', 'year'),
        <int64>json_get(value, 'date', 'month'),
        <int64>json_get(value, 'date', 'day'),
        0, 0, 0, 'UTC'
      ) else <datetime>{}
    ),
    precision := <date::DatePrecision>value['precision']
  )));
}
