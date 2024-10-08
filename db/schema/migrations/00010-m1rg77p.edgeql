CREATE MIGRATION m1rg77poeapfh3u77ttpitjsmh7ciw3d253a3hjkqdjdvqiblzv77a
    ONTO m1huuufuhetqpkemw73stzhyu6lujfm7abv7glzvvopwoefvi2ulnq
{
  ALTER TYPE people::PendingUserRequest {
      CREATE REQUIRED PROPERTY identity: tuple<first_name: std::str, last_name: std::str> {
          SET REQUIRED USING (<tuple<std::str, std::str>>{});
      };
  };
};
