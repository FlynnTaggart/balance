SET TIME ZONE 'Europe/Moscow';

-- Table spaces
ALTER TABLESPACE pg_global
    OWNER TO postgres;
ALTER TABLESPACE pg_default
    OWNER TO postgres;

-- Users
CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL NOT NULL,
    balance numeric NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
) TABLESPACE pg_default;

-- Services
CREATE TABLE IF NOT EXISTS services(
    id BIGSERIAL NOT NULL,
    name text NOT NULL,
    CONSTRAINT services_pkey PRIMARY KEY (id)
) TABLESPACE pg_default;

-- Reserves
CREATE TABLE IF NOT EXISTS reserves (
    order_id bigint NOT NULL,
    user_id bigint NOT NULL,
    service_id bigint NOT NULL,
    amount numeric NOT NULL,
    purchased bool NOT NULL,
    reserved_at timestamp,
    purchased_at timestamp,
    CONSTRAINT reserves_pkey PRIMARY KEY (order_id),
    CONSTRAINT fk_reserves_service FOREIGN KEY (service_id)
        REFERENCES services (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT fk_reserves_user FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
) TABLESPACE pg_default;

-- Indexes
CREATE INDEX reports ON reserves (service_id, amount, purchased_at) where purchased = true;