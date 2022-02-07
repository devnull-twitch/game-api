-- Add slot ID column

alter table character_items 
    add column slot_id int not null default 0;

---- create above / drop below ----

alter table character_items 
    drop column slot_id;
