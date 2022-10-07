-- Speichern 1 min, 5 min, 15 min, 1h, 3h , 4h, 6h

CREATE TYPE exchanges as ENUM(
    'ftx',
    'binance',
    'bybit',
    'deribit',
    'bitmex',
    'coinbase',
    'phemex'
);


CREATE TABLE IF NOT EXISTS index
(
    index_id BIGSERIAL PRIMARY KEY,
    name     varchar(64) NOT NULL unique
);


CREATE TABLE IF NOT EXISTS ticker
(
    ticker_id BIGSERIAL PRIMARY KEY,
    exchange  exchanges NOT NULL,
    ticker    varchar(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS ticker_index
(
    ticker_id     BIGINT references ticker (ticker_id) ON DELETE CASCADE NOT NULL,
    index_id      BIGINT references index (index_id) ON DELETE CASCADE NOT NULL,
    weight        int NOT NULL,
    excludevolume bool NOT NULL
);

CREATE TABLE IF NOT EXISTS ohclv
(
    index_id   BIGINT    NOT NULL,
    resolution int       NOT NULL,
    starttime  timestamp NOT NULL,
    open       float4    not null,
    high       float4    not null,
    close      float4    not null,
    low        float4    not null,
    volume     float4    not null,
    unique (index_id, starttime, resolution)
);

CREATE TABLE IF NOT EXISTS ftxlow
( -- < 1h
) INHERITS (ohclv);

CREATE TABLE IF NOT EXISTS ftxhigh
( -- > 1h
) INHERITS (ohclv);

CREATE TABLE IF NOT EXISTS binacelow
( -- < 1h
) INHERITS (ohclv);

CREATE TABLE IF NOT EXISTS binancehigh
( -- < 1h
) INHERITS (ohclv);

CREATE TABLE IF NOT EXISTS indexlow
( -- < 1h
) INHERITS (ohclv);

CREATE TABLE IF NOT EXISTS indexhigh
( -- < 1h
) INHERITS (ohclv);

CREATE TABLE IF NOT EXISTS deribitlow
( -- < 1h
) INHERITS (ohclv);

CREATE TABLE IF NOT EXISTS deribithigh
( -- < 1h
) INHERITS (ohclv);


---- create above / drop below ----

DROP TABLE IF EXISTS deribithigh;
DROP TABLE IF EXISTS deribitlow;
DROP TABLE IF EXISTS indexhigh;
DROP TABLE IF EXISTS indexlow;
DROP TABLE IF EXISTS binacelow;
DROP TABLE IF EXISTS binancehigh;
DROP TABLE IF EXISTS ftxhigh;
DROP TABLE IF EXISTS ftxlow;
DROP TABLE IF EXISTS ohclv;
DROP TABLE IF EXISTS ticker_index;
DROP TABLE IF EXISTS ticker;
DROP TABLE IF EXISTS  index;
DROP TYPE IF EXISTS exchanges;