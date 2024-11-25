module occurrence {

  type Identification extending default::Auditable {
    required taxon: taxonomy::Taxon;

    # Defines whether this occurrence is the first scientific description of
    # the taxon
    required is_type: bool {
      default := false;
    };

    identified_by: people::Person; # might be unknown

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



  scalar type QuantityType extending enum<Unknown, One, Several, Ten, Tens, Hundred>;

  type OccurrenceReport extending Occurrence {
    reported_by: people::Person;
    reference: references::Article;

    original_link: str; # link to original database
    in_collection: str;
    multi item_voucher: str;

    required quantity: QuantityType;
    content_description: str;

    multi link sequences := .<source_sample[is seq::ExternalSequence];
  }
}
