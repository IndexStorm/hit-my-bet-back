BEGIN;

DROP TABLE IF EXISTS prediction.markets;
DROP TYPE IF EXISTS prediction.market_resolution;
DROP TYPE IF EXISTS prediction.market_chain_status;

DROP SCHEMA IF EXISTS prediction;

COMMIT;
