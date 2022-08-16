# Fin Documents management app(in future more fin entities will be added)
### Tools:
- go 1.19
- postgres

### How to use
Run container with postgres:
```cmd
docker run -d --name fin-db -e POSTGRES_PASSWORD=postgres -v ${HOME}/pgdata/:/var/lib/postgresql/data -p 5432:5432 postgres
```
Migrate tables using migrate tool(https://github.com/golang-migrate/migrate):
```cmd
migrate -path ./schema -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" up
```
Example of table created:
```sql
create table docs (
 id           serial not null unique,
 type         varchar(100) not null,
 counterparty varchar(255) not null,
 amount       float8       not null,
 doc_currency varchar(10)  not null,
 amount_usd   float8       not null,
 doc_date     timestamp not null default now(),
 notes        varchar(4000) not null,
 created      timestamp not null default now(),
 updated      timestamp not null default now()
);
```

Building and running the application:
```cmd
go build -o app cmd/main.go
./app
```