CREATE MIGRATION m1t2ij46mbc7qowsixuauhblvwolammlhjvidhgjkeknj56epxmqyq
    ONTO m1je3ok6e6jrsqqnup6j6bvwc3d2jtv4tucbtoremqrmdle333jqja
{
  ALTER TYPE events::Spotting {
      CREATE PROPERTY comments: std::str;
  };
};
