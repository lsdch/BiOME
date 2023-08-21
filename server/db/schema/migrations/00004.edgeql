CREATE MIGRATION m15eeizjt2q7igd2sjkwy5oz6jicc5bmpggm7p4kfapn62bdblm4uq
    ONTO m1pjrvyeddw466xedtztu3rfjux6vjxy5q3rdgjvgfei6gbeor2gba
{
  ALTER TYPE people::Person {
      ALTER PROPERTY second_names {
          SET TYPE array<std::str> USING (<array<std::str>>[.second_names]);
      };
  };
};
