create table docs (
 id           serial not null unique,
 type         varchar(100) not null,
 counterparty varchar(255) not null,
 amount       float8       not null,
 doc_currency varchar(10)  not null,
 amount_usd   float8       not null,
 doc_date     timestamp not null default now(),
 notes        varchar(4000) not null,
 created_at   timestamp not null default now(),
 updated_at   timestamp not null default now()
);