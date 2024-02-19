CREATE MIGRATION m15hpcvcd2fpy7jcwn6kj36iu5jdw3ybxmo7f5t2jioungb4qzprlq
    ONTO m1xzxtbkorylvfsw2vrow4bnaaehnvnq2syat5goftchn5qgqnoxeq
{
  ALTER TYPE people::Person {
      CREATE PROPERTY role := (.user.role);
  };
};
