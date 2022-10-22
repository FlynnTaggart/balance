SET TIMEZONE="Europe/Moscow";

-- Table spaces
ALTER TABLESPACE pg_global
    OWNER TO postgres;
ALTER TABLESPACE pg_default
    OWNER TO postgres;

-- Users
CREATE TABLE IF NOT EXISTS users(
    id BIGSERIAL NOT NULL,
    balance numeric,
    CONSTRAINT users_pkey PRIMARY KEY (id)
) TABLESPACE pg_default;

-- Services
CREATE TABLE IF NOT EXISTS services(
    id BIGSERIAL NOT NULL,
    name text,
    CONSTRAINT services_pkey PRIMARY KEY (id)
) TABLESPACE pg_default;

-- Orders
CREATE TABLE IF NOT EXISTS orders(
    id BIGSERIAL NOT NULL,
    CONSTRAINT orders_pkey PRIMARY KEY (id)
) TABLESPACE pg_default;

-- Reserves
CREATE TABLE IF NOT EXISTS reserves (
    user_id bigint NOT NULL,
    service_id bigint NOT NULL,
    order_id bigint NOT NULL,
    amount numeric,
    reserved_at timestamp with time zone,
    CONSTRAINT reserves_pkey PRIMARY KEY (user_id, service_id, order_id),
    CONSTRAINT fk_reserves_order FOREIGN KEY (order_id)
        REFERENCES orders (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT fk_reserves_service FOREIGN KEY (service_id)
        REFERENCES services (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION,
    CONSTRAINT fk_reserves_user FOREIGN KEY (user_id)
        REFERENCES users (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
) TABLESPACE pg_default;

-- Reports
CREATE TABLE IF NOT EXISTS reports (
    id BIGSERIAL NOT NULL,
    service_id bigint,
    amount numeric,
    purchased_at timestamp with time zone,
    CONSTRAINT reports_pkey PRIMARY KEY (id),
    CONSTRAINT fk_reports_service FOREIGN KEY (service_id)
        REFERENCES services (id)
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
) TABLESPACE pg_default;