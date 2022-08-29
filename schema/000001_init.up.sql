create table users (
 id            serial not null unique,
 name          varchar(255) not null,
 surname       varchar(255) not null,
 email         varchar(255) not null unique,
 password      varchar(255) not null,
 registered_at timestamp not null default now()
);

create table docs (
 id           serial not null unique,
 type         varchar(100) not null,
 counterparty varchar(255) not null,
 amount       float8       not null,
 doc_currency varchar(10)  not null,
 amount_usd   float8       not null,
 doc_date     timestamp not null default now(),
 notes        varchar(4000) not null,
 author_id    int references users (id) on delete cascade not null,
 created_at   timestamp not null default now(),
 updater_id   int references users (id) on delete cascade not null,
 updated_at   timestamp not null default now()
);
