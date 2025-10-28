CREATE MIGRATION m1qoauopnuwhfpmsfzwrnhxrpzzwhbhh4tbpubrxqmuwsplu56vqrq
    ONTO m1tx3qdxao3zkirxehdavlzn6dux3ajqteyjmsqavnrmiftrkgf6iq
{
  ALTER TYPE references::Article {
      ALTER PROPERTY code {
          RESET default;
      };
  };
  DROP FUNCTION references::generate_article_code(authors: array<std::str>, year: std::int32);
  DROP GLOBAL references::alphabet;
};
