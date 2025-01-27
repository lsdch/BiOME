CREATE MIGRATION m1b5jq7vilsukd5c6wzl5dyhqn2mm3qtj3gynrlg7v2mi7vcpcltva
    ONTO m1jg732brzwr6p7x76eld7bs4nbc2wl2lgjbkb6spcw6pqjym543xa
{
  ALTER TYPE admin::InstanceSettings {
      ALTER PROPERTY name {
          RESET default;
      };
  };
};
