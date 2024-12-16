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
      annotation description := "Format like 'taxon_short_code[sampling_code]'";
      rewrite update using (
        if __specified__.code then .code # allow manual overwrite
        else .identification.taxon.code ++ "[" ++ .sampling.code ++ "]"
      );

      # 🚧 unsupported feature for now
      # > error: default expression cannot refer to links of inserted object
      # > this is a tempory implementation restriction
      # 🔗 see also: https://github.com/edgedb/edgedb/issues/7384
      # -----
      # default := gen_biomat_code(.identification, .sampling)
      # rewrite insert, update using (
      #   if re_test(r'^\d+$', .code)  then (
      #     .identification.taxon.code ++ "|" ++ events::event_code(.sampling.event)
      #   )
      #   else .code
      # );
    };
    index on (.code);

    code_history: array<tuple<code: str, time: datetime>> {
      readonly := true;
      rewrite update using (
        if __old__.code != .code then
        __old__.code_history ++ [(code := __old__.code, time := datetime_of_statement())]
        else .code_history
      );
    }

    # Defines whether this occurrence is the first scientific description of
    # the taxon
    required is_type: bool {
      default := false;
    };

    multi published_in: references::Article {
      original_source: bool;
    };

  };

  type InternalBioMat extending BioMaterial {
    # multi link content := .<biomat[is samples::Sample];
    multi link specimens := .<biomat[is samples::Specimen];
    multi link bundles := .<biomat[is samples::BundledSpecimens];
    multi link identified_taxa := (
      select distinct .specimens.identification.taxon ?? .identification.taxon
    );

    required is_homogenous := (
      select count(distinct .specimens.identification.taxon) <= 1
    );

    required is_congruent := (
      select assert_exists(
        .is_homogenous and (
          (not exists .specimens) or
          .identification.taxon in (assert_single(distinct .specimens.identification.taxon))
        )
      )
    );
  }

  scalar type QuantityType extending enum<Unknown, One, Several, Ten, Tens, Hundred>;

  type ExternalBioMat extending BioMaterial {
    original_link: str; # link to original database

    in_collection: str; # name of a collection where the specimen can be found
    multi item_vouchers: str; # specimen identifier(s) within the collection

    required quantity: QuantityType;
    content_description: str;

    required is_homogenous := (
      select count(distinct .sequences.identification.taxon) <= 1
    );

    required is_congruent := (
      select assert_exists(
        .is_homogenous and (
          (not exists .sequences) or
          .identification.taxon in (assert_single(distinct .sequences.identification.taxon))
        )
      )
    );

    multi link sequences := .<source_sample[is seq::ExternalSequence];
  }

  alias BioMaterialWithType := (
    select BioMaterial {
      *,
      required has_sequences := (
        exists ([is ExternalBioMat].sequences ?? [is InternalBioMat].specimens.sequences)
      ),
      required is_homogenous := [is ExternalBioMat].is_homogenous ?? [is InternalBioMat].is_homogenous ?? true,
      required is_congruent := [is ExternalBioMat].is_congruent ?? [is InternalBioMat].is_congruent ?? true,
      category := (
        if (BioMaterial is InternalBioMat) then "Internal"
        else if (BioMaterial is ExternalBioMat) then "External"
        else "Unknown"
      )
    }
  )
}
