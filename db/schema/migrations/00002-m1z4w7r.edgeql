CREATE MIGRATION m1z4w7raqyomlkgxtyi6vt25aoexfdhiyedzlqjhwzka7zrzml32pq
    ONTO m1uf7gqulw57cndzcygywqjen4geft6x3e7ckzsrfk6vhqzby46aaq
{
  ALTER TYPE location::Site {
      ALTER LINK country {
          RESET OPTIONALITY;
      };
  };
};
