create table if not exists people (
  username varchar(30) not null primary key,
  name varchar(50) not null unique,
  email varchar(100) not null unique,
  phone_number varchar(15),
  wfdf_number varchar(20) unique,
  origin_country varchar(30),

  created_at timestamp not null default now(),
  created_by varchar(50),
  updated_at timestamp not null default now(),
  updated_by varchar(50)
);
