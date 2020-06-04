create table "Event" (
  id serial primary key,
  title varchar not null,
  dateStart timestamptz not null,
  dateEnd timestamptz not null,
  description text,
  userId int not null
);