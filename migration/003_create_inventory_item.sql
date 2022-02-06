-- Character inventory data

create table character_items (
  id serial primary key,
  chracter_id int,
  item_id int not null,
  quantity int default 1,
  constraint fk_character foreign key(chracter_id) references characters(id)
);

---- create above / drop below ----

drop table character_items;
