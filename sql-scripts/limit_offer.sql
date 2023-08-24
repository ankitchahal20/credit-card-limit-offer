CREATE TABLE IF NOT EXISTS public.limit_offer
(
    id character varying COLLATE pg_catalog."default" NOT NULL,
    account_id character varying COLLATE pg_catalog."default" NOT NULL,
    limit_type character varying COLLATE pg_catalog."default",
    new_limit integer,
    offer_activation_time timestamp with time zone,
    offer_expiry_time timestamp with time zone,
    status character varying COLLATE pg_catalog."default",
    CONSTRAINT limit_offer_pkey PRIMARY KEY (id),
    CONSTRAINT limit_offer_account_id_fkey FOREIGN KEY (account_id)
        REFERENCES public.account (account_id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
)