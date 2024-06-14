CREATE MIGRATION m1mjwfnafilprirzrjnnxjys7doq7wilphuwbji7fwia5w2dgl7k2q
    ONTO m1g34f6pkqniznncg2piqjpoqm3flagsec2kuyryzemvyjoxd2m6la
{
  ALTER TYPE location::Country DROP EXTENDING default::Auditable;
};
