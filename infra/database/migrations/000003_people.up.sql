create table if not exists people (
  userName varchar(30) not null primary key,
  name varchar(50) not null unique,
  email varchar(100) not null unique,
  phoneNumber varchar(15),
  wfdfNumber varchar(20) unique,
  originCountry varchar(2),

  created_at timestamp not null default now(),
  created_by varchar(50),
  updated_at timestamp not null default now(),
  updated_by varchar(50)
);
