CREATE MIGRATION m15tsemmruo62xzzgb7jczincphvpnfbsvznardetpomlezd2ls2da
    ONTO m1atvhsx6ip3ds6b4aa5whl7hu4a66gnbp2lrrqcj6hlpd5ogncc4q
{
  ALTER SCALAR TYPE people::InstitutionKind EXTENDING enum<Lab, FundingAgency, SequencingPlatform, Other>;
};
