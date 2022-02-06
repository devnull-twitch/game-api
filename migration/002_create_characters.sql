-- Create character database

create table characters (
  id serial primary key,
  account_id int,
  character_name varchar not null,
  character_display varchar not null,
  base_color varchar,
  current_zone varchar,
  constraint fk_account foreign key(account_id) references accounts(id),
  constraint unique_char_name unique(character_name)
);

---- create above / drop below ----

drop table characters;