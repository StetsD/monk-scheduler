create table "Event" (
  id serial primary key,
  title varchar not null,
  dateStart timestamp with time zone not null,
  dateEnd timestamp with time zone not null,
  description text,
  userId int not null
);