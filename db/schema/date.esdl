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
}
