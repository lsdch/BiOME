CREATE MIGRATION m1c7j6p22ea67pi2rsiqq5tq5ye4wtbrhjppeakpxjjqajp7lsf5cq
    ONTO m1oywdlm6aub3uhvikw4qu34c2bblgfjxcw4tazbuckglznpay5irq
{
  ALTER TYPE datasets::Dataset {
      DROP LINK publication;
  };
  ALTER TYPE datasets::Dataset {
      CREATE MULTI LINK publications: references::Article {
          ON SOURCE DELETE ALLOW;
      };
  };
};
