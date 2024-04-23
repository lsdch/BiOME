CREATE MIGRATION m1egbksxywcdlfyuibcvplobjifj2mjdxztlg6rbqvh35ufesmpv3q
    ONTO m1x4b2i4zoipgyikyw3ja2mxfzibspsdzfj3v3fkbmlgm4i7tlvsla
{
  ALTER TYPE people::UserInvitation {
      CREATE REQUIRED LINK issued_by: people::User {
          SET default := (GLOBAL default::current_user);
      };
  };
};
