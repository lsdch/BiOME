CREATE MIGRATION m1l7trqua4xlbezhdxal6njx6wn6qovxgxt4emr2wrk5g2ebdqp7dq
    ONTO m1elskjwstzjcfust7yjko6jv53exahqsv5r5fuzoxrbyfvvu5ctga
{
  ALTER TYPE admin::Settings {
      DROP PROPERTY registration_enabled;
  };
};
