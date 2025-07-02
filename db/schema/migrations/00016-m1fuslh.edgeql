CREATE MIGRATION m1fuslhbl7c67z5afycwwiupsuix2c5k77aoxf2svyju64d25a5bsq
    ONTO m1gbrpj6kazku5fhety4tcujm7dxy7q7cnv2wvbnwqpuw5lehl65fq
{
  CREATE MODULE settings IF NOT EXISTS;
  CREATE TYPE settings::AbstractSettingsSpec EXTENDING default::Auditable {
      CREATE REQUIRED PROPERTY is_global: std::bool {
          SET default := false;
      };
      CREATE REQUIRED PROPERTY is_public: std::bool {
          SET default := false;
      };
      CREATE REQUIRED PROPERTY name: std::str {
          CREATE DELEGATED CONSTRAINT std::exclusive;
      };
      CREATE REQUIRED PROPERTY spec: std::json;
  };
  CREATE TYPE settings::DataFeedSpec EXTENDING settings::AbstractSettingsSpec;
  CREATE TYPE settings::MapToolPreset EXTENDING settings::AbstractSettingsSpec;
};
