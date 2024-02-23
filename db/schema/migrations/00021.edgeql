CREATE MIGRATION m166xwndlm53llpxfyylbufmpgnrgp3guvptnbqcfluyxwmxj4haba
    ONTO m1dyelrc2hzqdh2bzhodrmvtqgt5bfaaohngico7wajkkvspb4zh3a
{
  ALTER TYPE people::Person {
      CREATE PROPERTY comment: std::str;
  };
};
