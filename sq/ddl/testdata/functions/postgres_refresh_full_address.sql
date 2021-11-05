CREATE FUNCTION refresh_full_address() RETURNS trigger AS $$ BEGIN
    REFRESH MATERIALIZED VIEW CONCURRENTLY full_address;
    RETURN NULL;
END; $$ LANGUAGE plpgsql;
