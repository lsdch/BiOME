CREATE MIGRATION m1g34f6pkqniznncg2piqjpoqm3flagsec2kuyryzemvyjoxd2m6la
    ONTO m12wrtoe4mvegymkfndn4ugrl3n2dgsxhrugt2ctfavsvzmgdvf7vq
{
  ALTER TYPE location::Country {
      DROP LINK localities;
  };
};
