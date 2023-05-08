CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "order"
(
    order_uid          UUID primary key default uuid_generate_v4() ,
    track_number       varchar(255) not null unique,
    entry              varchar(255) not null,
    locale             varchar(65)  not null,
    internal_signature varchar(255),
    customer_id        varchar(255) not null,
    delivery_service   varchar(255) not null,
    shardkey           varchar(65)  not null,
    sm_id              int          not null,
    date_created       TIMESTAMP    not null   default CURRENT_TIMESTAMP,
    oof_shard          varchar(65)  not null
);

CREATE TABLE delivery
(
    order_uid UUID primary key,
    name      varchar(255) not null,
    phone     varchar(255) not null,
    zip       varchar(255) not null,
    city      varchar(255) not null,
    address   varchar(255) not null,
    region    varchar(255) not null,
    email     varchar(255) not null,

    FOREIGN KEY (order_uid) REFERENCES "order" (order_uid) ON DELETE CASCADE
);

CREATE TABLE payment
(
    order_uid     UUID primary key,
    transaction   varchar(255) not null,
    request_id    varchar(255) not null,
    currency      varchar(255) not null,
    provider      varchar(255) not null,
    amount        int          not null,
    payment_dt    varchar(255) not null,
    bank          varchar(255) not null,
    delivery_cost int          not null,
    goods_total   int          not null,
    custom_fee    int          not null,
    FOREIGN KEY (order_uid) REFERENCES "order" (order_uid) ON DELETE CASCADE
);

CREATE TABLE item
(
    id           serial primary key,
    order_uid    UUID         not null,
    track_number varchar(255) not null,
    price        int          not null,
    rid          varchar(255) not null,
    name         varchar(255) not null,
    sale         int          not null,
    size         varchar          not null default '0',
    total_price  int          not null,
    nm_id        int          not null,
    brand        varchar(255) not null,
    status       int          not null,

    FOREIGN KEY (order_uid) REFERENCES "order" (order_uid) ON DELETE CASCADE,
    FOREIGN KEY (track_number) REFERENCES "order" (track_number) ON DELETE CASCADE
);
