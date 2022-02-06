-- This is a sample migration.

create table accounts (
  id serial primary key,
  username varchar not null,
  password varchar,
  constraint unique_name unique(username)
);

---- create above / drop below ----

drop table accounts;
