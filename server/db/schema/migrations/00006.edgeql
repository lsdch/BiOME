CREATE MIGRATION m1lpnacdrieqhx25v66xaa5chojohpnne7btbdka5clfu4hmvgessq
    ONTO m1kghg74bjihfngeirongjeooiyy5vvcfb5qa6osb66gs7kjb6lejq
{
  ALTER TYPE people::User {
      ALTER PROPERTY name {
          RENAME TO login;
      };
  };
};
