CREATE MIGRATION m1d2lb2mgvfwvlhge5fxwv3l7ik2podkkplwius3nmnuker6owo7ca
    ONTO m1eyhpn7bbeew4vqoz2opjo3n77ipppshey3pc5hez2degqk42t6qq
{
  ALTER TYPE admin::Settings {
      CREATE PROPERTY geoapify_api_key: std::str;
  };
};
