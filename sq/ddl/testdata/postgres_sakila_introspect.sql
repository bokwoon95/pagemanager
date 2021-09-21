CREATE SCHEMA IF NOT EXISTS public;

CREATE EXTENSION IF NOT EXISTS btree_gist;

CREATE EXTENSION IF NOT EXISTS plpgsql;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE OR REPLACE FUNCTION public.last_update_trg()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$ BEGIN
    NEW.last_update = NOW();
    RETURN NEW;
END; $function$;

CREATE OR REPLACE FUNCTION public.refresh_full_address()
 RETURNS trigger
 LANGUAGE plpgsql
AS $function$ BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY full_address;
    RETURN NULL;
END; $function$;

CREATE TABLE IF NOT EXISTS public.actor (
    actor_id INTEGER NOT NULL GENERATED BY DEFAULT AS IDENTITY
    ,first_name TEXT NOT NULL
    ,last_name TEXT NOT NULL
    ,full_name TEXT GENERATED ALWAYS AS ((first_name || ' '::text) || last_name) STORED
    ,full_name_reversed TEXT GENERATED ALWAYS AS ((last_name || ' '::text) || first_name) STORED
    ,last_update TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP

    ,CONSTRAINT actor_actor_id_pkey PRIMARY KEY (actor_id)
);

CREATE TABLE IF NOT EXISTS public.address (
    address_id INTEGER NOT NULL GENERATED BY DEFAULT AS IDENTITY
    ,address TEXT NOT NULL
    ,address2 TEXT
    ,district TEXT NOT NULL
    ,city_id INTEGER NOT NULL
    ,postal_code TEXT
    ,phone TEXT NOT NULL
    ,last_update TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP

    ,CONSTRAINT address_address_id_pkey PRIMARY KEY (address_id)
);

CREATE TABLE IF NOT EXISTS public.category (
    category_id INTEGER NOT NULL GENERATED BY DEFAULT AS IDENTITY
    ,name TEXT NOT NULL
    ,last_update TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP

    ,CONSTRAINT category_category_id_pkey PRIMARY KEY (category_id)
);

CREATE TABLE IF NOT EXISTS public.city (
    city_id INTEGER NOT NULL GENERATED BY DEFAULT AS IDENTITY
    ,city TEXT NOT NULL
    ,country_id INTEGER NOT NULL
    ,last_update TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP

    ,CONSTRAINT city_city_id_pkey PRIMARY KEY (city_id)
);

CREATE TABLE IF NOT EXISTS public.country (
    country_id INTEGER NOT NULL GENERATED BY DEFAULT AS IDENTITY
    ,country TEXT NOT NULL
    ,last_update TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP

    ,CONSTRAINT country_country_id_pkey PRIMARY KEY (country_id)
);

CREATE TABLE IF NOT EXISTS public.customer (
    customer_id INTEGER NOT NULL GENERATED BY DEFAULT AS IDENTITY
    ,store_id INTEGER NOT NULL
    ,first_name TEXT NOT NULL
    ,last_name TEXT NOT NULL
    ,email TEXT
    ,address_id INTEGER NOT NULL
    ,active BOOLEAN NOT NULL DEFAULT true
    ,data JSONB
    ,create_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,last_update TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP

    ,CONSTRAINT customer_customer_id_pkey PRIMARY KEY (customer_id)
    ,CONSTRAINT customer_email_first_name_last_name_key UNIQUE (email, first_name, last_name)
    ,CONSTRAINT customer_email_key UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS public.film (
    film_id INTEGER NOT NULL GENERATED BY DEFAULT AS IDENTITY
    ,title TEXT NOT NULL
    ,description TEXT
    ,release_year INTEGER
    ,language_id INTEGER NOT NULL
    ,original_language_id INTEGER
    ,rental_duration INTEGER NOT NULL DEFAULT 3
    ,rental_rate NUMERIC(4,2) NOT NULL DEFAULT 4.99
    ,length INTEGER
    ,replacement_cost NUMERIC(5,2) NOT NULL DEFAULT 19.99
    ,rating TEXT DEFAULT 'G'::text
    ,special_features TEXT[]
    ,last_update TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,fulltext TSVECTOR

    ,CONSTRAINT film_film_id_pkey PRIMARY KEY (film_id)
    ,CONSTRAINT film_rating_check CHECK (rating = ANY (ARRAY['G'::text, 'PG'::text, 'PG-13'::text, 'R'::text, 'NC-17'::text]))
    ,CONSTRAINT film_release_year_check CHECK (release_year >= 1901 AND release_year <= 2155)
);

CREATE TABLE IF NOT EXISTS public.film_actor (
    film_id INTEGER NOT NULL
    ,actor_id INTEGER NOT NULL
    ,last_update TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.film_actor_review (
    film_id INTEGER NOT NULL
    ,actor_id INTEGER NOT NULL
    ,review_title TEXT NOT NULL DEFAULT ''::text COLLATE "C"
    ,review_body TEXT NOT NULL DEFAULT ''::text
    ,metadata JSONB
    ,last_update TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
    ,delete_date TIMESTAMP WITH TIME ZONE

    ,CONSTRAINT film_actor_review_check CHECK (length(review_body) > length(review_title))
    ,CONSTRAINT film_actor_review_film_id_actor_id_pkey PRIMARY KEY (film_id, actor_id)
);

CREATE TABLE IF NOT EXISTS public.film_category (
    film_id INTEGER NOT NULL
    ,category_id INTEGER NOT NULL
    ,last_update TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS public.inventory (
    inventory_id INTEGER NOT NULL GENERATED BY DEFAULT AS IDENTITY
    ,film_id INTEGER NOT NULL
    ,store_id INTEGER NOT NULL
    ,last_update TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP

    ,CONSTRAINT inventory_inventory_id_pkey PRIMARY KEY (inventory_id)
);

CREATE TABLE IF NOT EXISTS public.language (
    language_id INTEGER NOT NULL GENERATED BY DEFAULT AS IDENTITY
    ,name TEXT NOT NULL
    ,last_update TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP

    ,CONSTRAINT language_language_id_pkey PRIMARY KEY (language_id)
);

CREATE TABLE IF NOT EXISTS public.payment (
    payment_id INTEGER NOT NULL GENERATED BY DEFAULT AS IDENTITY
    ,customer_id INTEGER NOT NULL
    ,staff_id INTEGER NOT NULL
    ,rental_id INTEGER
    ,amount NUMERIC(5,2) NOT NULL
    ,payment_date TIMESTAMP WITH TIME ZONE NOT NULL

    ,CONSTRAINT payment_payment_id_pkey PRIMARY KEY (payment_id)
);

CREATE TABLE IF NOT EXISTS public.rental (
    rental_id INTEGER NOT NULL GENERATED BY DEFAULT AS IDENTITY
    ,rental_date TIMESTAMP WITH TIME ZONE NOT NULL
    ,inventory_id INTEGER NOT NULL
    ,customer_id INTEGER NOT NULL
    ,return_date TIMESTAMP WITH TIME ZONE
    ,staff_id INTEGER NOT NULL
    ,last_update TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP

    ,CONSTRAINT rental_range_excl EXCLUDE USING GIST (inventory_id WITH =, tstzrange(rental_date, return_date, '[]'::text) WITH &&)
    ,CONSTRAINT rental_rental_id_pkey PRIMARY KEY (rental_id)
);

CREATE TABLE IF NOT EXISTS public.staff (
    staff_id INTEGER NOT NULL GENERATED BY DEFAULT AS IDENTITY
    ,first_name TEXT NOT NULL
    ,last_name TEXT NOT NULL
    ,address_id INTEGER NOT NULL
    ,email TEXT
    ,store_id INTEGER
    ,active BOOLEAN NOT NULL DEFAULT true
    ,username TEXT NOT NULL
    ,password TEXT
    ,last_update TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
    ,picture BYTEA

    ,CONSTRAINT staff_staff_id_pkey PRIMARY KEY (staff_id)
);

CREATE TABLE IF NOT EXISTS public.store (
    store_id INTEGER NOT NULL GENERATED BY DEFAULT AS IDENTITY
    ,manager_staff_id INTEGER NOT NULL
    ,address_id INTEGER NOT NULL
    ,last_update TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP

    ,CONSTRAINT store_store_id_pkey PRIMARY KEY (store_id)
);

CREATE OR REPLACE VIEW public.actor_info AS  SELECT a.actor_id,
    a.first_name,
    a.last_name,
    jsonb_object_agg(c.name, ( SELECT jsonb_agg(f.title) AS jsonb_agg
           FROM film f
             JOIN film_category fc_1 ON fc_1.film_id = f.film_id
             JOIN film_actor fa_1 ON fa_1.film_id = f.film_id
          WHERE fc_1.category_id = c.category_id AND fa_1.actor_id = a.actor_id
          GROUP BY fa_1.actor_id)) AS film_info
   FROM actor a
     LEFT JOIN film_actor fa ON fa.actor_id = a.actor_id
     LEFT JOIN film_category fc ON fc.film_id = fa.film_id
     LEFT JOIN category c ON c.category_id = fc.category_id
  GROUP BY a.actor_id, a.first_name, a.last_name;

CREATE OR REPLACE VIEW public.customer_list AS  SELECT cu.customer_id AS id,
    (cu.first_name || ' '::text) || cu.last_name AS name,
    a.address,
    a.postal_code AS "zip code",
    a.phone,
    city.city,
    country.country,
        CASE
            WHEN cu.active THEN 'active'::text
            ELSE ''::text
        END AS notes,
    cu.store_id AS sid
   FROM customer cu
     JOIN address a ON a.address_id = cu.address_id
     JOIN city ON city.city_id = a.city_id
     JOIN country ON country.country_id = city.country_id;

CREATE OR REPLACE VIEW public.film_list AS  SELECT film.film_id AS fid,
    film.title,
    film.description,
    category.name AS category,
    film.rental_rate AS price,
    film.length,
    film.rating,
    jsonb_agg((actor.first_name || ' '::text) || actor.last_name) AS actors
   FROM category
     LEFT JOIN film_category ON film_category.category_id = category.category_id
     LEFT JOIN film ON film.film_id = film_category.film_id
     JOIN film_actor ON film_actor.film_id = film.film_id
     JOIN actor ON actor.actor_id = film_actor.actor_id
  GROUP BY film.film_id, film.title, film.description, category.name, film.rental_rate, film.length, film.rating;

CREATE MATERIALIZED VIEW IF NOT EXISTS public.full_address AS  SELECT country.country_id,
    city.city_id,
    address.address_id,
    country.country,
    city.city,
    address.address,
    address.address2,
    address.district,
    address.postal_code,
    address.phone,
    address.last_update
   FROM address
     JOIN city ON city.city_id = address.city_id
     JOIN country ON country.country_id = city.country_id;

CREATE OR REPLACE VIEW public.nicer_but_slower_film_list AS  SELECT film.film_id AS fid,
    film.title,
    film.description,
    category.name AS category,
    film.rental_rate AS price,
    film.length,
    film.rating,
    jsonb_agg((((upper("substring"(actor.first_name, 1, 1)) || lower("substring"(actor.first_name, 2))) || ' '::text) || upper("substring"(actor.last_name, 1, 1))) || lower("substring"(actor.last_name, 2))) AS actors
   FROM category
     LEFT JOIN film_category ON film_category.category_id = category.category_id
     LEFT JOIN film ON film.film_id = film_category.film_id
     JOIN film_actor ON film_actor.film_id = film.film_id
     JOIN actor ON actor.actor_id = film_actor.actor_id
  GROUP BY film.film_id, film.title, film.description, category.name, film.rental_rate, film.length, film.rating;

CREATE OR REPLACE VIEW public.sales_by_film_category AS  SELECT c.name AS category,
    sum(p.amount) AS total_sales
   FROM payment p
     JOIN rental r ON r.rental_id = p.rental_id
     JOIN inventory i ON i.inventory_id = r.inventory_id
     JOIN film f ON f.film_id = i.film_id
     JOIN film_category fc ON fc.film_id = f.film_id
     JOIN category c ON c.category_id = fc.category_id
  GROUP BY c.name
  ORDER BY (sum(p.amount)) DESC;

CREATE OR REPLACE VIEW public.sales_by_store AS  SELECT (ci.city || ','::text) || co.country AS store,
    (m.first_name || ' '::text) || m.last_name AS manager,
    sum(p.amount) AS total_sales
   FROM payment p
     JOIN rental r ON r.rental_id = p.rental_id
     JOIN inventory i ON i.inventory_id = r.inventory_id
     JOIN store s ON s.store_id = i.store_id
     JOIN address a ON a.address_id = s.address_id
     JOIN city ci ON ci.city_id = a.city_id
     JOIN country co ON co.country_id = ci.country_id
     JOIN staff m ON m.staff_id = s.manager_staff_id
  GROUP BY co.country, ci.city, s.store_id, m.first_name, m.last_name
  ORDER BY co.country, ci.city;

CREATE OR REPLACE VIEW public.staff_list AS  SELECT s.staff_id AS id,
    (s.first_name || ' '::text) || s.last_name AS name,
    a.address,
    a.postal_code AS "zip code",
    a.phone,
    ci.city,
    co.country,
    s.store_id AS sid
   FROM staff s
     JOIN address a ON a.address_id = s.address_id
     JOIN city ci ON ci.city_id = a.city_id
     JOIN country co ON co.country_id = ci.country_id;

CREATE INDEX IF NOT EXISTS actor_last_name_idx ON public.actor (last_name);

CREATE INDEX IF NOT EXISTS address_city_id_idx ON public.address (city_id);

CREATE INDEX IF NOT EXISTS city_country_id_idx ON public.city (country_id);

CREATE INDEX IF NOT EXISTS customer_address_id_idx ON public.customer (address_id);

CREATE INDEX IF NOT EXISTS customer_last_name_idx ON public.customer (last_name);

CREATE INDEX IF NOT EXISTS customer_store_id_idx ON public.customer (store_id);

CREATE INDEX IF NOT EXISTS film_fulltext_idx ON public.film USING GIST (fulltext);

CREATE INDEX IF NOT EXISTS film_language_id_idx ON public.film (language_id);

CREATE INDEX IF NOT EXISTS film_original_language_id_idx ON public.film (original_language_id);

CREATE INDEX IF NOT EXISTS film_title_idx ON public.film (title);

CREATE UNIQUE INDEX IF NOT EXISTS film_actor_actor_id_film_id_idx ON public.film_actor (actor_id, film_id);

CREATE INDEX IF NOT EXISTS film_actor_film_id_idx ON public.film_actor (film_id);

CREATE INDEX IF NOT EXISTS film_actor_review_misc ON public.film_actor_review (film_id, (substr(review_body, 2, 10)), (review_title || ' abcd'::text), ((metadata ->> 'score'::text)::integer)) INCLUDE (actor_id, last_update) WHERE delete_date IS NULL;

CREATE INDEX IF NOT EXISTS film_actor_review_review_body_idx ON public.film_actor_review (review_body);

CREATE INDEX IF NOT EXISTS film_actor_review_review_title_idx ON public.film_actor_review (review_title);

CREATE INDEX IF NOT EXISTS inventory_store_id_film_id_idx ON public.inventory (store_id, film_id);

CREATE INDEX IF NOT EXISTS payment_customer_id_idx ON public.payment (customer_id);

CREATE INDEX IF NOT EXISTS payment_staff_id_idx ON public.payment (staff_id);

CREATE INDEX IF NOT EXISTS rental_customer_id_idx ON public.rental (customer_id);

CREATE INDEX IF NOT EXISTS rental_inventory_id_idx ON public.rental (inventory_id);

CREATE UNIQUE INDEX IF NOT EXISTS rental_rental_date_inventory_id_customer_id_idx ON public.rental (rental_date, inventory_id, customer_id);

CREATE INDEX IF NOT EXISTS rental_staff_id_idx ON public.rental (staff_id);

CREATE UNIQUE INDEX IF NOT EXISTS store_manager_staff_id_idx ON public.store (manager_staff_id);

CREATE UNIQUE INDEX IF NOT EXISTS full_address_country_id_city_id_address_id_idx ON public.full_address (country_id, city_id, address_id) INCLUDE (country, city, address, address2);

CREATE TRIGGER actor_last_update_before_update_trg BEFORE UPDATE ON actor FOR EACH ROW EXECUTE FUNCTION last_update_trg();

CREATE TRIGGER address_refresh_full_address_trg AFTER INSERT OR DELETE OR UPDATE OR TRUNCATE ON address FOR EACH STATEMENT EXECUTE FUNCTION refresh_full_address();

CREATE TRIGGER city_last_update_before_update_trg BEFORE UPDATE ON address FOR EACH ROW EXECUTE FUNCTION last_update_trg();

CREATE TRIGGER category_last_update_before_update_trg BEFORE UPDATE ON category FOR EACH ROW EXECUTE FUNCTION last_update_trg();

CREATE TRIGGER city_last_update_before_update_trg BEFORE UPDATE ON city FOR EACH ROW EXECUTE FUNCTION last_update_trg();

CREATE TRIGGER city_refresh_full_address_trg AFTER INSERT OR DELETE OR UPDATE OR TRUNCATE ON city FOR EACH STATEMENT EXECUTE FUNCTION refresh_full_address();

CREATE TRIGGER country_last_update_before_update_trg BEFORE UPDATE ON country FOR EACH ROW EXECUTE FUNCTION last_update_trg();

CREATE TRIGGER country_refresh_full_address_trg AFTER INSERT OR DELETE OR UPDATE OR TRUNCATE ON country FOR EACH STATEMENT EXECUTE FUNCTION refresh_full_address();

CREATE TRIGGER customer_last_update_before_update_trg BEFORE UPDATE ON customer FOR EACH ROW EXECUTE FUNCTION last_update_trg();

CREATE TRIGGER film_fulltext_before_insert_update_trg BEFORE INSERT OR UPDATE ON film FOR EACH ROW EXECUTE FUNCTION tsvector_update_trigger('fulltext', 'pg_catalog.english', 'title', 'description');

CREATE TRIGGER film_last_update_before_update_trg BEFORE UPDATE ON film FOR EACH ROW EXECUTE FUNCTION last_update_trg();

CREATE TRIGGER film_actor_last_update_before_update_trg BEFORE UPDATE ON film_actor FOR EACH ROW EXECUTE FUNCTION last_update_trg();

CREATE TRIGGER film_actor_review_last_update_before_update_trg BEFORE UPDATE ON film_actor_review FOR EACH ROW EXECUTE FUNCTION last_update_trg();

CREATE TRIGGER film_category_last_update_before_update_trg BEFORE UPDATE ON film_category FOR EACH ROW EXECUTE FUNCTION last_update_trg();

CREATE TRIGGER inventory_last_update_before_update_trg BEFORE UPDATE ON inventory FOR EACH ROW EXECUTE FUNCTION last_update_trg();

CREATE TRIGGER language_last_update_before_update_trg BEFORE UPDATE ON language FOR EACH ROW EXECUTE FUNCTION last_update_trg();

CREATE TRIGGER rental_last_update_before_update_trg BEFORE UPDATE ON rental FOR EACH ROW EXECUTE FUNCTION last_update_trg();

CREATE TRIGGER staff_last_update_before_update_trg BEFORE UPDATE ON staff FOR EACH ROW EXECUTE FUNCTION last_update_trg();

CREATE TRIGGER store_last_update_before_update_trg BEFORE UPDATE ON store FOR EACH ROW EXECUTE FUNCTION last_update_trg();

ALTER TABLE IF EXISTS public.address
    ADD CONSTRAINT address_city_id_fkey FOREIGN KEY (city_id) REFERENCES public.city (city_id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE IF EXISTS public.city
    ADD CONSTRAINT city_country_id_fkey FOREIGN KEY (country_id) REFERENCES public.country (country_id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE IF EXISTS public.customer
    ADD CONSTRAINT customer_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.address (address_id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE IF EXISTS public.film
    ADD CONSTRAINT film_language_id_fkey FOREIGN KEY (language_id) REFERENCES public.language (language_id) ON UPDATE CASCADE ON DELETE RESTRICT
    ,ADD CONSTRAINT film_original_language_id_fkey FOREIGN KEY (original_language_id) REFERENCES public.language (language_id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE IF EXISTS public.film_actor
    ADD CONSTRAINT film_actor_actor_id_fkey FOREIGN KEY (actor_id) REFERENCES public.actor (actor_id) ON UPDATE CASCADE ON DELETE RESTRICT
    ,ADD CONSTRAINT film_actor_film_id_fkey FOREIGN KEY (film_id) REFERENCES public.film (film_id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE IF EXISTS public.film_actor_review
    ADD CONSTRAINT film_actor_review_film_id_actor_id_fkey FOREIGN KEY (film_id, actor_id) REFERENCES public.film_actor (film_id, actor_id) ON UPDATE CASCADE ON DELETE NO ACTION DEFERRABLE INITIALLY DEFERRED;

ALTER TABLE IF EXISTS public.film_category
    ADD CONSTRAINT film_category_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.category (category_id) ON UPDATE CASCADE ON DELETE RESTRICT
    ,ADD CONSTRAINT film_category_film_id_fkey FOREIGN KEY (film_id) REFERENCES public.film (film_id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE IF EXISTS public.inventory
    ADD CONSTRAINT inventory_film_id_fkey FOREIGN KEY (film_id) REFERENCES public.film (film_id) ON UPDATE CASCADE ON DELETE RESTRICT
    ,ADD CONSTRAINT inventory_store_id_fkey FOREIGN KEY (store_id) REFERENCES public.store (store_id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE IF EXISTS public.payment
    ADD CONSTRAINT payment_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customer (customer_id) ON UPDATE CASCADE ON DELETE RESTRICT
    ,ADD CONSTRAINT payment_rental_id_fkey FOREIGN KEY (rental_id) REFERENCES public.rental (rental_id) ON UPDATE CASCADE ON DELETE RESTRICT
    ,ADD CONSTRAINT payment_staff_id_fkey FOREIGN KEY (staff_id) REFERENCES public.staff (staff_id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE IF EXISTS public.rental
    ADD CONSTRAINT rental_customer_id_fkey FOREIGN KEY (customer_id) REFERENCES public.customer (customer_id) ON UPDATE CASCADE ON DELETE RESTRICT
    ,ADD CONSTRAINT rental_inventory_id_fkey FOREIGN KEY (inventory_id) REFERENCES public.inventory (inventory_id) ON UPDATE CASCADE ON DELETE RESTRICT
    ,ADD CONSTRAINT rental_staff_id_fkey FOREIGN KEY (staff_id) REFERENCES public.staff (staff_id) ON UPDATE CASCADE ON DELETE RESTRICT;

ALTER TABLE IF EXISTS public.staff
    ADD CONSTRAINT staff_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.address (address_id) ON UPDATE CASCADE ON DELETE RESTRICT
    ,ADD CONSTRAINT staff_store_id_fkey FOREIGN KEY (store_id) REFERENCES public.store (store_id) ON UPDATE NO ACTION ON DELETE NO ACTION;

ALTER TABLE IF EXISTS public.store
    ADD CONSTRAINT store_address_id_fkey FOREIGN KEY (address_id) REFERENCES public.address (address_id) ON UPDATE CASCADE ON DELETE RESTRICT
    ,ADD CONSTRAINT store_manager_staff_id_fkey FOREIGN KEY (manager_staff_id) REFERENCES public.staff (staff_id) ON UPDATE CASCADE ON DELETE RESTRICT;
