create table "user"
(
    username      text primary key,
    password_hash text               not null,
    is_system     bool default false not null
);

insert into "user" (username, password_hash, is_system)
values ('registration_gift', md5(random()::text), true),
       ('buy_item', md5(random()::text), true);

create table balance
(
    username text references "user" (username),
    amount   bigint not null,
    unique (username)
);

create table item
(
    name  text primary key,
    price bigint not null
);

insert into item (name, price)
values ('t-shirt', 80),
       ('cup', 20),
       ('book', 50),
       ('pen', 10),
       ('powerbank', 200),
       ('hoody', 300),
       ('umbrella', 200),
       ('socks', 10),
       ('wallet', 50),
       ('pink-hoody', 500);

create table inventory
(
    username  text references "user" (username),
    item_name text references item (name),
    quantity  bigint default 1 not null,
    unique (username, item_name)
);

create table transaction
(
    id      bigint generated always as identity primary key,
    "from"  text references "user" (username),
    "to"    text references "user" (username),
    amount  int                                 not null,
    created timestamp default current_timestamp not null
);

