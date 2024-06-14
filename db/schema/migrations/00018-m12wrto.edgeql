CREATE MIGRATION m12wrtoe4mvegymkfndn4ugrl3n2dgsxhrugt2ctfavsvzmgdvf7vq
    ONTO m1vev6tejtldvpak7nk7itzueju3bss34as6gy3zlckmvavjympfoq
{
  ALTER TYPE location::Locality DROP EXTENDING default::Auditable;
};
