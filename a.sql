create table usrs (
    id serial primary key,
    name varchar(16) not null unique,
    pass varchar(88) not null,

    check(length(name) >= 8)
);

create table ads (
    id serial primary key,
    title varchar(255) not null,
    content text not null, 
    image_url text,
    price int,
    user_id int references usrs (id) on delete cascade,

    check(length(title) >= 2),
    check(length(content) >= 2),
    check(price >= 1 and price <= 1000000)
);
