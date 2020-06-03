alter table "Event"
add column createdAt timestamp with time zone not null default CURRENT_TIMESTAMP,
    add column updateAt timestamp with time zone not null default CURRENT_TIMESTAMP,
    add column deleteAt timestamp with time zone;