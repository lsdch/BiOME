CREATE MIGRATION m1x4b2i4zoipgyikyw3ja2mxfzibspsdzfj3v3fkbmlgm4i7tlvsla
    ONTO m1gy2agbdoat3bsdqbe7apdaaw6n7poybyu7fo72zta4lxy4imigvq
{
  ALTER TYPE admin::Settings {
      CREATE REQUIRED PROPERTY registration_enabled := (SELECT
          (EXISTS (.email) AND .instance.allow_contributor_signup)
      );
  };
};
