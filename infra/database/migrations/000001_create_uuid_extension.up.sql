-- this migration can't be ran automatically in prod, because the RDS user does not have proper permission to do it
create extension if not exists "uuid-ossp";
