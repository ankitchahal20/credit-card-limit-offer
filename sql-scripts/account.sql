CREATE TABLE IF NOT EXISTS public.account
(
    account_id character varying COLLATE pg_catalog."default" NOT NULL,
    customer_id character varying COLLATE pg_catalog."default" NOT NULL,
    account_limit integer,
    per_transaction_limit integer,
    last_account_limit integer,
    last_per_transaction_limit integer,
    account_limit_update_time timestamp with time zone,
    per_transaction_limit_update_time timestamp with time zone,
    CONSTRAINT account_pkey PRIMARY KEY (account_id)
)