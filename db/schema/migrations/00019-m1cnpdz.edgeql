CREATE MIGRATION m1cnpdzskteshpateawedfj34e2i5zcs22qgoqt5dqvim3kzcead4a
    ONTO m12eobx6drp6r2akldqqg24s5wadmhrnxgrtiztdj3lu2uhrkbknvq
{
  CREATE SCALAR TYPE traits::Category EXTENDING enum<Morphology, Physiology, Ecology, Behaviour, LifeHistory, HabitatPref>;
  CREATE SCALAR TYPE traits::TraitDefinitionScope EXTENDING enum<Specimen, Taxon>;
  CREATE ABSTRACT TYPE traits::AbstractTrait {
      CREATE REQUIRED PROPERTY category: traits::Category;
      CREATE REQUIRED PROPERTY name: std::str;
      CREATE CONSTRAINT std::exclusive ON ((.category, .name));
      CREATE REQUIRED MULTI PROPERTY scopes: traits::TraitDefinitionScope {
          CREATE CONSTRAINT std::exclusive;
      };
      CREATE PROPERTY description: std::str;
  };
  ALTER TYPE traits::QualitativeTrait {
      CREATE PROPERTY category: traits::Category {
          SET REQUIRED USING (<traits::Category>{});
      };
  };
  ALTER TYPE traits::QualitativeTrait {
      DROP PROPERTY modalities;
  };
  ALTER TYPE traits::QualitativeTrait {
      CREATE PROPERTY scopes: traits::TraitDefinitionScope {
          SET REQUIRED USING (<traits::TraitDefinitionScope>{});
      };
  };
  ALTER TYPE traits::QualitativeTrait {
      DROP PROPERTY scopes;
  };
  ALTER TYPE traits::QualitativeTrait {
      CREATE PROPERTY scopes: traits::TraitDefinitionScope {
          SET REQUIRED USING (<traits::TraitDefinitionScope>{});
      };
      CREATE REQUIRED PROPERTY value: std::str {
          SET REQUIRED USING (<std::str>{});
      };
      EXTENDING traits::AbstractTrait,
      default::Auditable LAST;
      ALTER PROPERTY name {
          RESET OPTIONALITY;
          DROP OWNED;
          RESET TYPE;
      };
      CREATE CONSTRAINT std::exclusive ON ((.name, .value));
      ALTER PROPERTY category {
          RESET OPTIONALITY;
          DROP OWNED;
          RESET TYPE;
      };
  };
  ALTER TYPE traits::QualitativeTrait {
      ALTER PROPERTY scopes {
          RESET OPTIONALITY;
          DROP OWNED;
          RESET TYPE;
      };
  };
  DROP TYPE traits::QualitativeTraitTag;
};
