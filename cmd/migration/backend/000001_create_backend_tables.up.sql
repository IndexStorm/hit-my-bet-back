BEGIN;

CREATE SCHEMA prediction;

CREATE TYPE prediction.market_chain_status AS ENUM (
  'PENDING',
  'NEED_RETRY',
  'CONFIRMED'
  );

CREATE TYPE prediction.market_resolution AS ENUM (
  'UNRESOLVED',
  'TIE',
  'YES',
  'NO'
  );

CREATE TABLE prediction.markets
(
  id              TEXT                           NOT NULL,
  chain_status    prediction.market_chain_status NOT NULL,
  title           TEXT                           NOT NULL,
  description     TEXT,
  creator_pubkey  TEXT                           NOT NULL,
  resolver_pubkey TEXT                           NOT NULL,
  market_pubkey   TEXT                           NOT NULL,
  resolution      prediction.market_resolution   NOT NULL,
  created_at      pg_catalog.timestamptz         NOT NULL,
  open_through    pg_catalog.timestamptz         NOT NULL,
  PRIMARY KEY (id)
);

COMMIT;
