CREATE MIGRATION m17jfq4b5kfsl5ttztbmjd4diur55rszloxmsuoa3us3fycms5pnxa
    ONTO m14tdpsnqdetfb2lid2rygk2xuk7kjvgzgjir3nmuyehntad4ywgza
{
  ALTER TYPE people::User {
      ALTER LINK identity {
          SET REQUIRED USING (<people::Person>{});
      };
      ALTER PROPERTY is_active {
          SET default := true;
          RESET EXPRESSION;
          RESET CARDINALITY;
          RESET OPTIONALITY;
          SET TYPE std::bool;
      };
  };
  ALTER TYPE people::PendingUserRequest {
      CREATE REQUIRED PROPERTY email_verified: std::bool {
          SET default := false;
      };
      ALTER PROPERTY motive {
          RESET OPTIONALITY;
      };
  };
  ALTER TYPE people::User {
      DROP PROPERTY email_confirmed;
  };
  ALTER TYPE people::UserInvitation {
      ALTER PROPERTY dest {
          RENAME TO email;
      };
  };
};
