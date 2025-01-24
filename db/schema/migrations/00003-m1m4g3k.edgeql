CREATE MIGRATION m1jg732brzwr6p7x76eld7bs4nbc2wl2lgjbkb6spcw6pqjym543xa
    ONTO m1lvc4nif6vlxn23m7rtth7dpexhmsbqmk6g2gcsixuqobh7qstpta
{
  ALTER TYPE seq::Sequence {
      EXTENDING default::CodeIdentifier LAST;
  };
};
