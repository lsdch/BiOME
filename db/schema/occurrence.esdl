module occurrence {

  type Identification extending default::Auditable {
    required taxon: taxonomy::Taxon;

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

    comments: str;
  }

  abstract type BioMaterial extending Occurrence {
    required code : str {
      constraint exclusive;
      annotation description := "Format like 'taxon_short_code|sampling_code'";
      # rewrite insert, update using ((
      #   .identification.taxon.code ++ "|" ++ events::event_code(.sampling.event)
      # ));
    };
    index on (.code);

    # Defines whether this occurrence is the first scientific description of
    # the taxon
    required is_type: bool {
      default := false;
    };

    multi published_in: references::Article;
  };

  type InternalBioMat extending BioMaterial {
    multi link content := .<biomat[is samples::Sample];
    multi link specimens := .<biomat[is samples::Specimen];
    multi link bundles := .<biomat[is samples::BundledSpecimens];
    multi link identified_taxa := (
      select distinct .specimens.identification.taxon ?? .identification.taxon
    );
  }

  scalar type QuantityType extending enum<Unknown, One, Several, Ten, Tens, Hundred>;

  type ExternalBioMat extending BioMaterial {
    original_link: str; # link to original database

    in_collection: str; # name of a collection where the specimen can be found
    multi item_vouchers: str; # specimen identifier(s) within the collection

    required quantity: QuantityType;
    content_description: str;

    multi link sequences := .<source_sample[is seq::ExternalSequence];
  }

  alias BioMaterialWithType := (
    select BioMaterial {
      *,
      published_in: { * },
      sampling: { * },
      identification: { * },
      meta: { * },
      type := (
        if (BioMaterial is InternalBioMat) then "Internal"
        else if (BioMaterial is ExternalBioMat) then "External"
        else "Unknown"
      )
    }
  )
}
