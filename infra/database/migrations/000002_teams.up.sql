create table if not exists teams (
  slug varchar(30) not null primary key,
  name varchar(50) not null unique,
  description text,
  origin_country varchar(30) not null,

  created_at timestamp not null default now(),
  created_by varchar(50),
  updated_at timestamp not null default now(),
  updated_by varchar(50)
);
