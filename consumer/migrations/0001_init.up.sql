create table if not exists wallets (
    id bigserial primary key
);

create table if not exists currency_accounts (
    id bigserial primary key,
    currency_code varchar(8) not null,
    amount double precision not null check(Amount >= 0),
    wallet_id bigint not null,
    unique(wallet_id, currency_code),
    foreign key (wallet_id) references wallets (id)
);

create table if not exists transactions (
    id bigserial primary key,
    amount double precision not null check(Amount != 0),
    currency_code varchar(8) not null,
    status varchar(16) not null, -- Created | Error | Success
    wallet_id bigint not null,
    foreign key (wallet_id) references wallets (id)
);

insert into wallets (id) values (1);
insert into currency_accounts (currency_code, amount, wallet_id) values ('RUB', 0, 1);
insert into currency_accounts (currency_code, amount, wallet_id) values ('USDT', 0, 1);
insert into currency_accounts (currency_code, amount, wallet_id) values ('EUR', 0, 1);