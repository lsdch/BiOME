CREATE MIGRATION m1jyqyxtr7irk7lpxqwwz6tv7qw5rg35qq3k5aiodtk3biwtvpoiwq
    ONTO m1ez5o5finxnvighgosksewvxafspvlquayl4zj776ve45edviklpq
{
  ALTER TYPE default::Meta {
      ALTER ANNOTATION std::title := 'Tracking data modifications';
      CREATE PROPERTY lastUpdated := ((.modified ?? .created));
  };
};
