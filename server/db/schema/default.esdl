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

module default {
  global current_user_id: uuid;

  abstract type Auditable {
    annotation title := "Auto-generation of timestamps";

    created: datetime {
      rewrite insert using (datetime_of_statement())
    }
    modified: datetime {
      rewrite update using (datetime_of_statement())
    }
  }

  abstract type Vocabulary {
    annotation title := "An extensible list of terms";

    required label: str { constraint exclusive };
    required code: str {
      annotation title := "An expressive, unique, user-generated uppercase alphanumeric code";
      constraint exclusive;
      constraint min_len_value(2);
      constraint max_len_value(8);
    };
    description: str;
  }

  type Conservation extending Vocabulary,Auditable;

  type Picture {
    legend: str;
    required path: str;
  }
}

module reference {
  type Article extending default::Auditable {
    required authors: str;
    required year: int16 {
      constraint min_value(1500)
    };
    required title: str;
    comments: str;
   }
}


module taxonomy {

  scalar type Rank extending enum<Kingdom, Phylum, Class, Order, Family, Genus, Species, Subspecies>;
  # Ignore FORM, VARIETY, UNRANKED

  scalar type TaxonStatus extending enum<Accepted, Synonym, Unclassified>;

  type Taxon extending default::Auditable {
    required GBIF_ID: int32 {
      constraint exclusive
    };
    required name: str;
    required rank: Rank;

    required status: TaxonStatus;
    required code: str {
      constraint exclusive;
      rewrite insert, update using (
        with chopped := str_split(.name, " "),
        suffix := "[syn]" if .status = TaxonStatus.Synonym else ""
        select (.name ++ suffix)
        if not .rank in {Rank.Species, Rank.Subspecies}
        else str_upper(chopped[0][:3]) ++ array_join(chopped[1:], "_")
      )
    };
    required anchor: bool {
      default := false
    }
    authorship: str;

    parent: Taxon;
    constraint exclusive on ((.name, .status));
    constraint expression on (exists .parent) except (.rank = Rank.Kingdom);
    # required is_classified: bool; # is present in established taxonomy
    # cannot be validated or synonym if not classified
  }

  # type Taxon extending Taxon {
  #   overloaded required code: str {
  #     constraint exclusive;
  #     rewrite insert, update using (
  #       with suffix := "[syn]" if .status = TaxonStatus.Synonym else ""
  #       select (.name ++ suffix)
  #     );
  #   };
  #   parent: Taxon;
  #    # bracket_authorship ?
  #   multi link children := .<parent[is Taxon];
  # }

  # type Species extending Taxon {
  #   required parent: Taxon;
  #   overloaded required code: str {
  #     constraint exclusive;
  #     rewrite insert, update using (
  #       with chopped := str_split(.name, " ")
  #       select(str_upper(chopped[0][:3]) ++ array_join(chopped[1:], "_"))
  #     )
  #   };
  #   overloaded required authorship: str;
  # }
}

module location {
  type Country extending default::Auditable {
    annotation description := "Countries as defined in the ISO 3166-1 norm.";
    required name: str {
      constraint exclusive;
      readonly := true;
    };
    required code: str {
      constraint exclusive;
      readonly := true;
      constraint min_len_value(2);
      constraint max_len_value(2);
    };
    multi link localities := .<country[is Locality]
  }

  type Locality extending default::Auditable {
    region: str;
    municipality: str;
    constraint exclusive on ((.region, .municipality));

    required code: str {
      constraint exclusive; # generate using municipality|region[country_code]
    };

    required country: Country;
    multi link sites := .<locality[is Site]
  }

  type Habitat extending default::Vocabulary, default::Auditable;
  type AccessPoint extending default::Vocabulary, default::Auditable;


  scalar type CoordinateMaxPrecision extending enum<"10m", "100m", "1Km", "10Km", "100Km", "Unknown">;

  type Site extending default::Auditable {
    required name : str { constraint exclusive };
    required code : str {
      annotation title := "Site identifier";
      annotation description := "A short, unique, user-generated, alphanumeric identifier. Recommended size is 8.";
      constraint exclusive;
      constraint min_len_value(4);
      constraint max_len_value(8);
    }
    description: str;

    required habitat: Habitat {
      annotation title := "The type of habitat that is the target of the sampling."
    };
    required access_point: AccessPoint {
      annotation title := "The actual point where the sampling is performed.";
      annotation description := "Some habitats may not be directly accessible, and sampling may have to be done on a location that acts as a proxy for the target habitat.";
    };

    required locality: Locality;

    required coordinates: tuple<
      precision: CoordinateMaxPrecision,
      latitude: float32,
      longitude: float32
    > {
      constraint expression on (
        (.latitude <= 90 and .latitude >= -90 and .longitude <= 180 and .longitude >= -180)
      );
    };

    altitudeRange: tuple<min:float32, max:float32> {
      annotation title := "The site elevation in meters"
    };

    multi link abiotic_measurements := .<site[is event::AbioticMeasurement];

    multi link samplings := .<site[is event::Sampling];
  }
}

module event {

  type Program extending default::Auditable {
    required label: str { constraint exclusive };
    required code: str { constraint exclusive };
    required multi managers: people::Person;
    multi funding_agencies: people::Institution;
    start_year: int16 {
      constraint min_value(1900);
    };
    end_year: int16;
    constraint expression on (.start_year <= .end_year);
    # TODO : constraint end year > start year + change to integer
    comments: str;
  }

  type SamplingMethod extending default::Vocabulary, default::Auditable;

  type PlannedSampling extending default::Auditable {
    required site: location::Site;
    required multi found_by: people::Person;
    multi targetTaxa: taxonomy::Taxon; # TODO: required ?
    comments: str;
  }

  abstract type Event extending default::Auditable {
    required site: location::Site;
    required performed_on: tuple<date: datetime, precision: date::DatePrecision> {
      constraint date::required_unless_unknown(__subject__.date, __subject__.precision);
      rewrite insert, update using (
        (
          date := date::rewrite_date(__subject__.performed_on),
          precision := __subject__.performed_on.precision
        )
      )
    };
    multi programs: Program;
  }


  type AbioticParameter extending default::Vocabulary, default::Auditable;

  type AbioticMeasurement extending event::Event {
    required param: AbioticParameter;
    required value: float32;
    constraint exclusive on ((.site, .param, .performed_on))
  }

  scalar type SamplingTarget extending enum<Community, Unknown, Taxa>;

  type Sampling extending Event {
    property generated_code := (
      select .site.code ++
      "_" ++ <str>datetime_get(.performed_on.date, 'year') ++
      <str>datetime_get(.performed_on.date, 'month')
    );
    required code : str {
      annotation title := "Unique sampling identifier, auto-generated at sampling creation.";
      constraint exclusive;
      rewrite insert using (
        (select .generated_code ++ "_" ++ <str>(select count(Sampling) filter Sampling.code = __subject__.generated_code)
          if (select exists Sampling filter Sampling.code = __subject__.generated_code)
          else .generated_code)
      )
      # TODO : generate as SITE_YEARMONTH_ORDER
    };

    # WARNING : during migration, remove pseudo-field when no sampling was performed
    multi methods: SamplingMethod;
    multi fixatives: default::Conservation;

    required sampling_target: SamplingTarget;
    multi target_taxa: taxonomy::Taxon {
      annotation title := "Taxonomic groups that were the target of the sampling effort"
    };
    # handle in application logic
    # constraint expression on (exists .target_taxa) except (.sampling_target != SamplingTarget.Taxa);

    sampling_duration: duration {
      annotation title := "Sampling duration measured in minutes."
    };

    required is_donation: bool;

    comments: str;
    multi link biomaterials := .<sampling[is samples::BioMaterial]
  }
}

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
    required sampling: event::Sampling;
    required identification: Identification {
      constraint exclusive
    };
    comments: str;
  }

  type Dataset extending default::Auditable {
    required name: str;
    required multi events: event::Sampling;
    required assembled_by: people::Person;
    origin: str;
    multi published_in: reference::Article;
    comments: str;
  }

  scalar type QuantityType extending enum<Exact, Unknown, One, Several, Ten, Tens, Hundred>;

  type OccurrenceReport extending Occurrence {
    recorded_by: people::Person;
    reference: reference::Article;

    original_link: str; # link to original database

    specimen_available : tuple<collection: str, item: str> {
      constraint exclusive
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
  }
}


module storage {

  abstract type Storage extending default::Auditable {
    required label: str {
      constraint exclusive
    };
    required code: str {
      constraint exclusive
    }
    description: str;
    required collection: Collection;
  }

  type BioMatStorage extending Storage;
  type SlideStorage extending Storage;
  type DNAStorage extending Storage;

  type Collection {
    required label: str {constraint exclusive };
    required code: str { constraint exclusive };
    required taxon: taxonomy::Taxon;
    comments: str;

    multi link bio_mat_storages := .<collection[is BioMatStorage];
    multi link slide_storages := .<collection[is SlideStorage];
    multi link DNA_storages := .<collection[is DNAStorage];
  }
}

module samples {

  type BioMaterial extending occurrence::Occurrence {
    # Think about splitting between internal and external
    required code : str {
      constraint exclusive;
      rewrite insert, update using (
        .identification.taxon.code ++ "|" ++ .sampling.code
      )
    };
    required created_on: tuple<date: datetime, precision: date::DatePrecision> {
      constraint date::required_unless_unknown(.date, .precision);
      rewrite insert, update using (
        (date := date::rewrite_date(.created_on), precision :=.created_on.precision)
      )
    };
    # required date: default::DateWithPrecision;
    required multi sorted_by: people::Person;
    multi published_in: reference::Article;
  }

  type ContentType extending default::Vocabulary, default::Auditable;

  abstract type Tube {
    required biomat: BioMaterial;

    required code: str {
      constraint min_len_value(1);
      constraint max_len_value(10);
    }; # TODO: add length constraints
    constraint exclusive on ((.biomat, .code));

    required content_type: ContentType;
    required conservation: default::Conservation;
    comments: str;
  }

  type MultipleTube extending Tube, default::Auditable {
    required specimen_quantity: int16 {
      constraint min_value(2)
    };
  }

  type IndividualTube extending Tube, default::Auditable {
    required specimen: Specimen;
  }

  type ChelexTube extending Tube, default::Auditable {
    required number: int16;
    required color: str;
    required origin: IndividualTube;
  }


  # type SpecimenMolID extending default::Auditable {
  #   required specimen: Specimen;
  #   required sequence: seq::AssembledSequence;
  #   required identified_by: people::Person;
  # }

  type Specimen extending default::Auditable {
    required morphological_code: str; #auto-generated
    curated_by: people::Person;

    molecular_code: str { constraint exclusive }; # auto-generated
    required molecular_number: str {
      rewrite insert using (
        <str>(select count(Specimen) filter Specimen.tube.biomat.sampling.site = __subject__.tube.biomat.sampling.site)
      )
    };
    # TODO : generate sequence at the level of site
    # is string because retrocompatibility

    multi pictures: default::Picture;
    comments: str;
    multi link slides := .<specimen[is Slide];
    link tube := .<specimen[is IndividualTube];
  }

  type Slide extending default::Auditable {
    required specimen: Specimen;
    required label: str; # the physical label
    property code: str {
      constraint exclusive;
      rewrite insert, update using (
        (select array_join([.storage.collection.code, .storage.code, <str>.storage_position], "_"))
      )
    }; # as {collectionCode}_{containerCode}_{slideNumberInBox}
    required created_on: tuple<date: datetime, precision: date::DatePrecision> {
      constraint date::required_unless_unknown(.date, .precision);
      rewrite insert, update using (
        (date := date::rewrite_date(.created_on), precision :=.created_on.precision)
      )
    };
    # date: default::DateWithPrecision;
    required multi mounted_by: people::Person;

    required storage: storage::SlideStorage;
    required storage_position: int16;
    constraint exclusive on ((.storage, .storage_position));

    multi pictures: default::Picture;
    comment: str;
  }
}



module seq {

  type BatchRequest extending default::Auditable {
    # TODO : need to link back from DNA tube
    required label: str;
    required target_gene: Gene; # generally 16S at first
    required multi samples: samples::ChelexTube;
    required requested_on: datetime;
    required requested_by: people::Person;
    multi requested_to: people::Person; # or group of persons ?
    required is_ready: bool {
      default := false
    };
    achieved_on: datetime;
    comments: str;
  }

  scalar type DNAQuality extending enum<Unknown, Contaminated, Bad, Good>;



  type DNAExtractionMethod extending default::Vocabulary, default::Auditable;

  type DNA extending default::Auditable {
    required specimen: samples::Specimen;
    required code: str { constraint exclusive }; # maybe autogenerate, talk with Lara
    required concentration: float32 {
      annotation title := "DNA concentration in ng/µL";
      constraint min_value(1e-6);
    };
    required multi extracted_by: people::Person;
    required extraction_method: DNAExtractionMethod;
    required extracted_on: tuple<date: datetime, precision: date::DatePrecision> {
      constraint date::required_unless_unknown(.date, .precision);
      rewrite insert, update using (
        (date := date::rewrite_date(.extracted_on), precision :=.extracted_on.precision)
      )
    };
    # required extractionDate: default::DateWithPrecision;
    required quality : DNAQuality;
    required is_empty: bool {
      default := false
    };
    required stored_in: storage::DNAStorage;
    comments: str;
   }


  type Gene extending default::Auditable {
    required name: str { constraint exclusive }; # e.g. cytochrome oxydase
    required code: str { constraint exclusive }; # COI
    description: str;
  }

  abstract type Primer extending default::Auditable {
    required label: str { constraint exclusive };
    required code: str { constraint exclusive };
    description: str;
  }

  type PCRForwardPrimer extending Primer, default::Auditable;

  type PCRReversePrimer extending Primer, default::Auditable;

  type PCRSpecificity extending default::Auditable {
    required label: str { constraint exclusive };
    description: str;
  }

  scalar type PCRQuality extending enum<Failure, Acceptable, Good, Unknown>;

  type PCR extending default::Auditable {
    # Store files ?
    required DNA: DNA;
    required gene: Gene;
    required code: str { constraint exclusive };
    required forward_primer : PCRForwardPrimer;
    required reverse_primer : PCRReversePrimer;
    required quality: PCRQuality;
    required performed_on: tuple<date: datetime, precision: date::DatePrecision> {
      constraint date::required_unless_unknown(.date, .precision);
      rewrite insert, update using (
        (date := date::rewrite_date(__subject__.performed_on), precision := __subject__.performed_on.precision)
      )
    };
    # required date: default::DateWithPrecision;
    comments: str;
   }


  type ChromatoPrimer extending Primer, default::Auditable;

  scalar type ChromatoQuality extending enum<Contaminated, Failure, Ok, Unknown>;


  type SequencingInstitute extending default::Auditable {
    required name: str { constraint exclusive };
    comments: str;
  }

  type Chromatogram extending default::Auditable {
    # store chromatogram files ?
    required PCR: PCR;
    required code: str { constraint exclusive };
    required YAS_number: str {constraint exclusive};
    required primer: ChromatoPrimer;
    required quality: ChromatoQuality;
    required provider: SequencingInstitute;
    comments: str;
  }

  abstract type Sequence {
    required code: str { constraint exclusive };
    required alignmentCode: str { constraint exclusive };
    required taxon: taxonomy::Taxon;
    comments: str;
  }

  type AssembledSequence extending Sequence, default::Auditable {

    accessionNumber: str { constraint exclusive };

    required multi assembled_by: people::Person;
    multi published_in: reference::Article;
    required multi chromatograms: Chromatogram;
  }

  type ExternalSequence extending Sequence, default::Auditable {
    # Talk with Christophe about data migration
    required sampling: event::Sampling;

    required accessionNumber: str { constraint exclusive };
    required specimenIdentifier: str;

    primaryTaxon: taxonomy::Taxon;
    required gene: Gene;
  }
}

module datasets {

  type Alignment extending default::Auditable {
    required label: str;
    required multi sequences: seq::Sequence;
    comments: str;
  }


  type MOTUDataset extending default::Auditable {
    required label: str;
    multi link MOTUs := .<dataset[is MOTU];
    published_in: reference::Article;
    required multi generated_by: people::Person;
    comments: str;
  }

  type DelimitationMethod extending default::Vocabulary, default::Auditable;

  type MOTU extending default::Auditable {
    required dataset: MOTUDataset;
    required number: int16;
    required method: DelimitationMethod;
    required multi sequences: seq::Sequence;
  }
}

module people {

  type Institution {
    required name: str { constraint exclusive };
    comments: str;
  }

  type Person {
    required first_name: str;
    second_names: str; # TODO check spelling ; up to 3 values
    required last_name: str;

    property full_name := .first_name ++ ' ' ++ .last_name;

    contact: str;
    multi institution: Institution;
  }


  scalar type UserRole extending enum<Guest, Contributor, ProjectMember, Admin>;

  type User {
    required name: str {constraint exclusive }; #autogenerated as first letter of first name + last name
    required email: str { constraint exclusive };
    required password: str;
    required verified: bool {
      default := false
    };

    role: UserRole;
    required identity: Person;
  }
}



module traits {
  scalar type Category extending enum<Morphology, Physiology, Ecology, Behaviour, LifeHistory, HabitatPref>;

  abstract type AbstractTrait {
    required category: Category;
    required name: str;
    description: str;
    constraint exclusive on ((.category, .name));
  }


  type QuantitativeTrait extending AbstractTrait;
  type QualitativeTrait extending AbstractTrait;

  type FuzzyTrait extending AbstractTrait {
    required multi modalities: tuple<name: str, range: tuple<min:int16, max:int16>>;
  }

  abstract type QualitativeMeasurement {
    required trait: QualitativeTrait;
    required value: str;
  }

  abstract type QuantitativeMeasurement {
    required trait: QuantitativeTrait;
    required value: float32;
  }

  abstract type SpeciesTrait {
    required species: taxonomy::Taxon;
    reference: reference::Article;
    expert_opinion: people::Person;
    constraint expression on (exists .reference or exists .expert_opinion);
  }



  type QualitativeIndividualTrait extending QualitativeMeasurement {
    required method: str;
  }

  type QuantitativeIndividualTrait extending QuantitativeMeasurement {
    required method: str;
  }

  type QualitativeSpeciesTrait extending QuantitativeMeasurement, SpeciesTrait;

  type QuantitativeSpeciesTrait extending QualitativeMeasurement, SpeciesTrait;



  type FuzzyTraitValue extending SpeciesTrait {
    required trait: FuzzyTrait;
    required multi values: tuple<name: str, value: int16> {
      annotation description := "Must be validated in the application logic."
    };
  }
}