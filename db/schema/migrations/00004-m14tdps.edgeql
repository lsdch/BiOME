CREATE MIGRATION m14tdpsnqdetfb2lid2rygk2xuk7kjvgzgjir3nmuyehntad4ywgza
    ONTO m1xvrlz4my3tzflk3k6vf2bnb2e2k5zzzljxslvmgpyjsder5s6aaq
{
  ALTER TYPE people::PendingUserRequest {
      DROP LINK user;
  };
  ALTER TYPE people::PendingUserRequest {
      CREATE REQUIRED PROPERTY email: std::str {
          SET REQUIRED USING (<std::str>{});
      };
  };
};
