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
  global current_user := (
    select people::User filter .id = global current_user_id limit 1
  );
  abstract annotation example;

  type Meta {
    annotation title := "Tracking of data changes";

    required created: datetime {
      default := datetime_of_statement();
    };
    # required created_by: people::Person {
    #   default := (select people::Person filter .user.id = global current_user_id);
    # };

    modified: datetime {
      rewrite update using (datetime_of_statement())
    };
    # modified_by: people::Person {
    #   rewrite update using (
    #     select people::Person filter .user.id = global current_user_id
    #   )
    # };
  }

  abstract type Auditable {
    annotation title := "Auto-generation of timestamps";

    required meta: Meta {
      constraint exclusive;
      on source delete delete target;
      on target delete restrict;
      default := (insert Meta {});
      rewrite update using (update .meta set { created := .created });
    }

    # required created: datetime {
    #   default := datetime_of_statement()
    # }
    # modified: datetime {
    #   rewrite update using (datetime_of_statement())
    # }
  }

  abstract type Vocabulary {
    annotation title := "An extensible list of terms";

    required label: str {
      constraint exclusive;
      constraint min_len_value(4);
    };
    required code: str {
      annotation title := "An expressive, unique, user-generated uppercase alphanumeric code";
      constraint exclusive;
      constraint min_len_value(2);
      constraint max_len_value(8);
    };
    description: str;
  }

  type Conservation extending Vocabulary, Auditable {
    annotation description := "Describes a conservation method for a sample."
  };

  type Picture {
    legend: str;
    required path: str {
      constraint exclusive;
    };
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
      constraint exclusive;
    };
    required name: str {
      constraint min_len_value(4);
    };

    property slug := (str_replace(.name, " ", "-"));

    required rank: Rank;

    required status: TaxonStatus;
    required code: str {
      # constraint exclusive;
      rewrite insert, update using (
        with chopped := str_split(.name, " "),
        suffix := "[syn]" if .status = TaxonStatus.Synonym else ""
        select (.name ++ suffix)
        if not .rank in {Rank.Species, Rank.Subspecies}
        else str_upper(chopped[0][:3]) ++ array_join(chopped[1:], "_") ++ suffix
      )
    };
    required anchor: bool {
      annotation description := "Signals whether this taxon was manually imported";
      default := false;
    }
    authorship: str;

    optional parent: Taxon {
      on target delete delete source;
    };
    constraint expression on (exists .parent) except (.rank = Rank.Kingdom);

    constraint exclusive on ((.name, .status));

    multi link children := .<parent[is Taxon];
  }
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
      constraint exclusive;
      annotation description := "Format like 'municipality|region[country_code]'";
      rewrite insert, update using (
        .municipality ?? "NA" ++ "|"
        ++ .region ?? "NA" ++ "|"
        ++ .country.code
      )
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
      annotation title := "The type of habitat that is the target of the sampling.";
      on target delete restrict;
    };
    required access_point: AccessPoint {
      annotation title := "The actual point where the sampling is performed.";
      annotation description := "Some habitats may not be directly accessible, and sampling may have to be done on a location that acts as a proxy for the target habitat.";
      on target delete restrict;
    };

    required locality: Locality {
      on target delete restrict;
    };

    required coordinates: tuple<
      precision: CoordinateMaxPrecision,
      latitude: float32,
      longitude: float32
    > {
      constraint expression on (
        (.latitude <= 90 and .latitude >= -90
        and .longitude <= 180 and .longitude >= -180)
      );
    };

    altitudeRange: tuple<min:float32, max:float32> {
      annotation title := "The site elevation in meters";
    };


    multi link abiotic_measurements := .<site[is event::AbioticMeasurement];
    multi link samplings := .<site[is event::Sampling];
    multi link spottings := .<site[is event::Spotting];
  }

  type SiteDataset extending default::Auditable {
    required label: str {
      constraint min_len_value(4);
      constraint max_len_value(40);
    }
    description: str;
    multi sites: Site;
    required multi maintainers: people::Person;
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
    comments: str;
  }

  abstract type Event extending default::Auditable {
    required site: location::Site;

    required multi performed_by: people::Person;
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

  type Spotting extending Event {
    multi target_taxa: taxonomy::Taxon;
  }

  type AbioticParameter extending default::Vocabulary, default::Auditable;

  type AbioticMeasurement extending event::Event {
    required param: AbioticParameter;
    required value: float32;
    related_sampling: Sampling;
    constraint exclusive on ((.site, .param, .performed_on))
  }

  scalar type SamplingTarget extending enum<Community, Unknown, Taxa>;

  type SamplingMethod extending default::Vocabulary, default::Auditable;

  type Sampling extending Event {
    property generated_code := (
      select .site.code ++
      "_" ++ <str>datetime_get(.performed_on.date, 'year') ++
      <str>datetime_get(.performed_on.date, 'month')
    );
    required code : str {
      annotation title := "Unique sampling identifier, auto-generated at sampling creation.";
      annotation description := "Format : SITE_YEARMONTH_NUMBER. The NUMBER suffix is not appended if the site and month tuple is unique.";
      annotation default::example := "SOMESITE_202301 ; SOMESITE_202301_1";

      constraint exclusive;
      rewrite insert using (
        (select
          .generated_code ++ "_" ++
          <str>(select count(Sampling) filter Sampling.code = __subject__.generated_code
        ) if (select exists Sampling filter Sampling.code = __subject__.generated_code)
          else .generated_code)
      )
    };

    # WARNING : during migration, remove pseudo-field when no sampling was performed
    multi methods: SamplingMethod;
    multi fixatives: default::Conservation; # TODO : what is the relation with conservation of samples down the stream ?

    required sampling_target: SamplingTarget;
    multi target_taxa: taxonomy::Taxon {
      annotation title := "Taxonomic groups that were the target of the sampling effort"
      # handle in application logic
      # constraint expression on (exists .target_taxa) except (.sampling_target != SamplingTarget.Taxa);
    };

    sampling_duration: duration;

    required is_donation: bool;

    comments: str;

    multi link measurements := .<related_sampling[is AbioticMeasurement];
    multi link samples := .<sampling[is samples::BioMaterial];
    multi link reports := .<sampling[is occurrence::OccurrenceReport];
    multi link external_seqs := .<sampling[is seq::ExternalSequence];

    multi link occurring_taxa := (
      with ext_samples_no_seqs := (select .reports filter not exists .sequences)
      select distinct
      (ext_samples_no_seqs.identification.taxon) union
      (.external_seqs.identification.taxon) union
      (.samples.identified_taxa)
    );
  }

  type EventDataset extending default::Auditable {
    required name: str {
      constraint min_len_value(4);
      constraint max_len_value(40);
    };

    multi samplings: Sampling {
      on target delete allow;
    };
    multi abiotic_measurements: AbioticMeasurement {
      on target delete allow;
    };
    multi spottings: Spotting {
      on target delete allow;
    };

    required multi maintainers: people::Person {
      on target delete restrict;
    };
    multi published_in: reference::Article {
      on target delete allow;
    };
    comments: str;
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


module storage {

  abstract type Storage extending default::Auditable {
    required label: str {
      constraint exclusive;
      constraint min_len_value(4);
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
    required label: str {
      constraint exclusive;
      constraint min_len_value(4);
    };
    required code: str {
      constraint exclusive;
      constraint min_len_value(4)
    };
    required taxon: taxonomy::Taxon;
    required maintainers: people::Person;
    comments: str;

    multi link bio_mat_storages := .<collection[is BioMatStorage];
    multi link slide_storages := .<collection[is SlideStorage];
    multi link DNA_storages := .<collection[is DNAStorage];
  }
}

module samples {

  type BioMaterial extending occurrence::Occurrence {
    # TODO IMPORTANT :
    # Maybe move down the identification on the level of tube
    # this would be more flexible, and it is always possible to suggest splitting
    # the biomat bundle in the UI when identifications dont concur across tubes.

    # Think about splitting between internal and external
    required code : str {
      constraint exclusive;
      annotation description := "Format like 'taxon_code|sampling_code'";
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
    required multi sorted_by: people::Person;
    multi published_in: reference::Article;
    multi link content := .<biomat[is Sample];
    multi link specimens := .<biomat[is Specimen];
    multi link bundles := .<biomat[is BundledSpecimens];
    multi link identified_taxa := (select distinct .specimens.identification.taxon ?? .identification.taxon);
  }

  type ContentType extending default::Vocabulary, default::Auditable;

  abstract type Sample extending default::Auditable {
    required biomat: BioMaterial;
    required type: ContentType;
    required conservation: default::Conservation;

    required number: int16 {
      annotation description := "Incremental number that discriminates between tubes having the same type in a bio material lot. Used to generate the tube code.";
      rewrite insert using (
        select 1 + (
          max(Sample.number
            filter Sample.biomat = __subject__.biomat
            and Sample.type = __subject__.type
          ) ?? 0)
      )
    };
    property tube := (.type.code ++ <str>.number);
    # required tube: tuple<number: int16, code: str> {
    #   rewrite insert using (
    #     with biomat := __subject__.biomat,
    #     type := __subject__.type,
    #     number := (select 1 + (
    #       max(Sample.tube.number
    #         filter Sample.biomat = biomat
    #         and Sample.type = type
    #       ) ?? 0))
    #     select (number := number, code := .type.code ++ <str>number)
    #   );
    # }
    comments: str;
  }

  type BundledSpecimens extending Sample {
    annotation description := "A tube containing several specimens.";
    required quantity: int16 {
      constraint min_value(2)
    };
  }

  type Specimen extending Sample {
    annotation description := "A single specimen isolated in a tube.";
    required morphological_code: str {
      constraint exclusive;
      rewrite insert, update using (
        (__subject__.biomat.code ++ "[" ++ .tube ++ "]")
      )
    };

    dissected_by: people::Person;

    molecular_code: str {
      constraint exclusive;
    }; # auto-generated
    required molecular_number: str {
      rewrite insert using (
        <str>(select count(Specimen) filter Specimen.biomat.sampling.site = __subject__.biomat.sampling.site)
      );
      # is string because retrocompatibility
    };

    multi link sequences := .<specimen[is seq::AssembledSequence];
    link identification := (
      (select .sequences filter .is_reference).identification
    );

    multi pictures: default::Picture {
      on source delete delete target;
      # APP: remember to also delete files
    };
    multi link slides := .<specimen[is Slide];

  }

  type Slide extending default::Auditable {
    required specimen: Specimen;
    required label: str; # the physical label
    property code: str {
      annotation description := "Generated as '{collectionCode}_{containerCode}_{slidePositionInBox}'";
      constraint exclusive;
      rewrite insert, update using (
        (select array_join([.storage.collection.code, .storage.code, <str>.storage_position], "_"))
      )
    };

    required created_on: tuple<date: datetime, precision: date::DatePrecision> {
      constraint date::required_unless_unknown(.date, .precision);
      rewrite insert, update using (
        (date := date::rewrite_date(.created_on), precision :=.created_on.precision)
      )
    };

    required multi mounted_by: people::Person;

    required storage: storage::SlideStorage;
    required storage_position: int16;
    constraint exclusive on ((.storage, .storage_position));

    multi pictures: default::Picture {
      on source delete delete target;
    };
    comment: str;
  }
}



module seq {

  type BatchRequest extending default::Auditable {
    required label: str;
    required target_gene: Gene;

    requested_on: datetime {
      annotation description := "If empty, the request is a draft and can not be processed yet."
    };
    required requested_by: people::Person;

    multi requested_to: people::Person; # or group of persons ?
    achieved_on: datetime;

    comments: str;
    required multi content: DNA;
  }

  scalar type DNAQuality extending enum<Unknown, Contaminated, Bad, Good>;

  type DNAExtractionMethod extending default::Vocabulary, default::Auditable;

  type DNA extending default::Auditable {
    required chelex_tube: tuple<color:str, number:int16>;
    required specimen: samples::Specimen;
    required code: str { constraint exclusive };
    # maybe autogenerate, talk with Lara
    # maybe use the container as a referential ?

    concentration: float32 {
      annotation title := "DNA concentration in ng/µL";
      constraint min_value(1e-3);
    };

    required multi extracted_by: people::Person;
    required extraction_method: DNAExtractionMethod;
    required extracted_on: tuple<date: datetime, precision: date::DatePrecision> {
      constraint date::required_unless_unknown(.date, .precision);
      rewrite insert, update using (
        (date := date::rewrite_date(.extracted_on), precision :=.extracted_on.precision)
      )
    };

    required quality : DNAQuality;
    required is_empty: bool {
      default := false
    };
    required stored_in: storage::DNAStorage;
    multi link PCRs := .<DNA[is PCR];
    comments: str;
  }


  type Gene extending default::Auditable {
    required name: str {
      annotation default::example := "cytochrome oxydase";
      constraint exclusive;
    };
    required code: str {
      annotation default::example := "COI";
      constraint exclusive;
    };
    description: str;
    required motu: bool; # TODO : discuss with Florian
  }

  abstract type Primer extending default::Auditable {
    required label: str { constraint exclusive };
    required code: str { constraint exclusive };
    sequence: str {
      constraint min_len_value(5);
    };
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
    comments: str;
    multi link chromatograms := .<PCR[is Chromatogram];

    # Questions for Lara :
    # - change quality to be a 0-3 scale, with 0 being a failure, then 1,2 or 3 stars
    # - Nested PCR: do we want the initial PCR to be saved ?
    #     or only the nested one, while keeping the information on the initial primers ?
    # - dilution: should we store the dilution factor as well ?
    # - Utiliser un système de tags combinables pour la spécificité ?
    # - how useful would storing files be ?

    # Notes spécificité :
    # - 3KB: taille de fragment suffisant pour avoir à la fois le 16S et le COI (?)
   }


  type ChromatoPrimer extending Primer, default::Auditable;

  scalar type ChromatoQuality extending enum<Contaminated, Failure, Ok, Unknown>;


  type SequencingInstitute extending default::Auditable {
    required name: str { constraint exclusive };
    comments: str;
  }

  type Chromatogram extending default::Auditable {
    # store chromatogram files ? talk with hosting service about space availability
    required PCR: PCR;
    required YAS_number: str { constraint exclusive };
    required primer: ChromatoPrimer;
    required code: str {
      constraint exclusive;
      rewrite insert using (
        .code ?? (select .YAS_number ++ "|" ++ .primer.code)
      );
    };
    required quality: ChromatoQuality;
    # What is the difference between 'NON DISPONIBLE' and "INCONNU" ?

    required provider: SequencingInstitute;
    comments: str;
    multi link sequences := .<chromatograms[is AssembledSequence]
  }

  abstract type Sequence {
    required code: str { constraint exclusive };
    required gene: Gene;
    comments: str;

    legacy: tuple<id: int32, code: str, alignment_code: str> {
      annotation description := "Legacy identifiers for retrocompatibility with data stored in GOTIT.";
    };
  }

  type AssembledSequence extending Sequence, default::Auditable {
    required alignmentCode: str { constraint exclusive };
    accession_number: str {
      annotation description := "The NCBI accession number, if the sequence was uploaded.";
      constraint exclusive;
    };
    required identification: occurrence::Identification;

    multi published_in: reference::Article;
    required multi assembled_by: people::Person;
    required multi chromatograms: Chromatogram;
    required specimen: samples::Specimen;
    required is_reference: bool {
      annotation description := "Whether this sequence should be used as the reference for the identification of the specimen";
    };
    constraint exclusive on ((.specimen, .is_reference)) except (not .is_reference);
  }


  scalar type ExternalSeqType extending enum<NCBI, PERSCOM, LEGACY> {
    annotation description := "The sequence origin. 'PERSCOM' is 'Personal communication', 'Legacy' indicates the sequence originates from the lab but could not be registered as such due to missing required metadata.";
  };

  type ExternalSequence extending Sequence, occurrence::Occurrence, default::Auditable {

    required type: ExternalSeqType;
    source_sample: occurrence::OccurrenceReport;

    accession_number: str {
      annotation description := "NCBI accession number of the original sequence.";
      constraint exclusive;
      constraint min_len_value(6);
      constraint max_len_value(12);
    };
    constraint expression on (exists .accession_number) except (.type != ExternalSeqType.NCBI);

    required specimen_identifier: str {
      annotation description := "An identifier for the organism from which the sequence was produced, provided in the original source";
    };
    original_taxon: str {
      annotation description := "The verbatim identification provided in the original source.";
    };
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
    second_names: array<str>; # TODO check spelling ; up to 3 values
    required last_name: str;

    property full_name := .first_name ++ ' ' ++ .last_name;

    contact: str;
    multi institution: Institution;
  }


  scalar type UserRole extending enum<Guest, Contributor, ProjectMember, Admin>;

  type User {
    # Maybe inherit from Person ?
    required name: str {constraint exclusive }; #autogenerated as first letter of first name + last name
    required email: str { constraint exclusive };
    required password: str;
    required verified: bool {
      default := false
    };

    role: UserRole;
    required identity: Person {
      constraint exclusive;
      on target delete restrict;
      on source delete delete target;
    };
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