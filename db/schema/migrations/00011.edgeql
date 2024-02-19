CREATE MIGRATION m1atvhsx6ip3ds6b4aa5whl7hu4a66gnbp2lrrqcj6hlpd5ogncc4q
    ONTO m1bzpahvfvejjec2edpihodpfdzkj4bswrnlhr2x42sgsw4f3apypa
{
  CREATE SCALAR TYPE people::InstitutionKind EXTENDING enum<Lab, FoundingAgency, SequencingPlatform, Other>;
  ALTER TYPE people::Institution {
      CREATE REQUIRED PROPERTY kind: people::InstitutionKind {
          SET REQUIRED USING (<people::InstitutionKind>{'Other'});
      };
  };
};
