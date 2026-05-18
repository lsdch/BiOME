CREATE MIGRATION m1ao2w4urmaz4y6za2ubwo25argca4kbyxf5327iltkonuvrageyla
    ONTO m1macbp2mbzhhbgblmehk4t2d6rnniwntgcproyqrpzfulawleg2eq
{
  ALTER TYPE location::Site {
      DROP PROPERTY last_visited;
  };
};
