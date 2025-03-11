module traits {

  scalar type Category extending enum<Morphology, Physiology, Ecology, Behaviour, LifeHistory, HabitatPref>;
  scalar type TraitDefinitionScope extending enum<Specimen, Taxon>;

  abstract type AbstractTrait {
    required category: Category;
    required name: str;
    description: str;
    required multi scopes: TraitDefinitionScope {
      constraint exclusive;
    };
    constraint exclusive on ((.category, .name));
  }

  type QualitativeTrait extending AbstractTrait, default::Auditable {
    required value: str;
    constraint exclusive on ((.name, .value));
  }
}

# module traits {
#   scalar type Category extending enum<Morphology, Physiology, Ecology, Behaviour, LifeHistory, HabitatPref>;

#   abstract type AbstractTrait {
#     required category: Category;
#     required name: str;
#     description: str;
#     constraint exclusive on ((.category, .name));
#   }


#   type QuantitativeTrait extending AbstractTrait;
#   type QualitativeTrait extending AbstractTrait;

#   type FuzzyTrait extending AbstractTrait {
#     required multi modalities: tuple<name: str, range: tuple<min:int16, max:int16>>;
#   }

#   abstract type QualitativeMeasurement {
#     required trait: QualitativeTrait;
#     required value: str;
#   }

#   abstract type QuantitativeMeasurement {
#     required trait: QuantitativeTrait;
#     required value: float32;
#   }

#   abstract type SpeciesTrait {
#     required species: taxonomy::Taxon;
#     reference: references::Article;
#     expert_opinion: people::Person;
#     constraint expression on (exists .reference or exists .expert_opinion);
#   }



#   type QualitativeIndividualTrait extending QualitativeMeasurement {
#     required method: str;
#   }

#   type QuantitativeIndividualTrait extending QuantitativeMeasurement {
#     required method: str;
#   }

#   type QualitativeSpeciesTrait extending QuantitativeMeasurement, SpeciesTrait;

#   type QuantitativeSpeciesTrait extending QualitativeMeasurement, SpeciesTrait;



#   type FuzzyTraitValue extending SpeciesTrait {
#     required trait: FuzzyTrait;
#     required multi values: tuple<name: str, value: int16> {
#       annotation description := "Must be validated in the application logic."
#     };
#   }
# }