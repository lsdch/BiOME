module occurrence {

  type Identification extending default::Auditable {
    required taxon: taxonomy::Taxon;
    required identified_by: people::Person;
    required identified_on: tuple<date: datetime, precision: date::DatePrecision>{
      constraint date::required_unless_unknown(__subject__.date, __subject__.precision);
      rewrite insert, update using (
        (date := date::rewrite_date(.identified_on), precision :=.identified_on.precision)
      )
    };
  }

  abstract type Occurrence extending default::Auditable {
    required sampling: events::Sampling;
    required identification: Identification {
      constraint exclusive;
      on source delete delete target;
    };
    # required multi link identifications := (select .identification);
    comments: str;
  }



  scalar type QuantityType extending enum<Exact, Unknown, One, Several, Ten, Tens, Hundred>;

  type OccurrenceReport extending Occurrence {
    reported_by: people::Person;
    reference: reference::Article;

    original_link: str; # link to original database

    specimen_available : tuple<collection: str, item: str> {
      constraint exclusive;
    };

    required quantity : tuple<precision: QuantityType, exact: int16> {
      constraint expression on (
        (.precision = QuantityType.Exact and .exact > 0) or
        (.precision != QuantityType.Exact)
      );
      rewrite insert, update using (
        ((precision := __subject__.quantity.precision, exact := -1))
        if __subject__.quantity.precision != QuantityType.Exact
        else __subject__.quantity
      )
    };
    multi link sequences := .<source_sample[is seq::ExternalSequence];
  }
}
