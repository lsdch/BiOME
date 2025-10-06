CREATE MIGRATION m1e3oa3eutnj6gn7meietebiytmsilypesd4hsr66md3bh36oa2cna
    ONTO m1z3yrk4amxwng2hqa2nlkg2ibx6juobzz577appsckeo2bjlp2xdq
{
  CREATE REQUIRED GLOBAL location::SITE_SEARCH_THRESHOLD -> std::float32 {
      SET default := 0.7;
  };
};
