alter table "Event"
add column createdAt timestamptz  not null default CURRENT_TIMESTAMP,
    add column updateAt timestamptz  not null default CURRENT_TIMESTAMP,
    add column deleteAt timestamptz;