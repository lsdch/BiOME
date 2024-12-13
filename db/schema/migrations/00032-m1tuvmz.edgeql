CREATE MIGRATION m1tuvmztmirbyofvrlqw5bv63q7dvf6gf3yj3xmnxzc22eme4jcswq
    ONTO m1zuqwlcbfunmetdcaotdj3erkxxogxfy7ksyoryxduu6q4lemitzq
{
  ALTER TYPE seq::Sequence EXTENDING default::Auditable LAST;
};
