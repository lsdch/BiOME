CREATE MIGRATION m1k5kyyo6zy73u6dtpeiirmzqqeoylvsthn4szy4ryddr4q3mvjx5q
    ONTO m1qoyv3egwolhsu4zfrcgo5igpbgsu637zjxijuqerd4q3w6sh5hgq
{
  ALTER TYPE people::User {
      DROP PROPERTY email_public;
  };
};
