CREATE TABLE IF NOT EXISTS pm_site (
    site_id UUID
    ,domain TEXT
    ,subdomain TEXT

    ,CONSTRAINT pm_site_site_id_pkey PRIMARY KEY (site_id)
    ,CONSTRAINT pm_site_domain_subdomain_key UNIQUE (domain, subdomain)
);

CREATE TABLE IF NOT EXISTS pm_url (
    site_id UUID
    ,path TEXT
    ,plugin TEXT
    ,params JSON

    ,CONSTRAINT pm_url_site_id_path_pkey PRIMARY KEY (site_id, path)
    ,CONSTRAINT pm_url_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id)
);

CREATE INDEX IF NOT EXISTS pm_url_site_id_idx ON pm_url (site_id);

CREATE TABLE IF NOT EXISTS pm_template_data (
    site_id UUID
    ,langcode TEXT
    ,data_file TEXT
    ,data JSON

    ,CONSTRAINT pm_template_data_site_id_langcode_data_file_pkey PRIMARY KEY (site_id, langcode, data_file)
    ,CONSTRAINT pm_template_data_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id)
);

CREATE INDEX IF NOT EXISTS pm_template_data_site_id_idx ON pm_template_data (site_id);

CREATE TABLE IF NOT EXISTS pm_user (
    user_id UUID
    ,username TEXT
    ,email TEXT
    ,name TEXT
    ,password_hash TEXT

    ,CONSTRAINT pm_user_user_id_pkey PRIMARY KEY (user_id)
    -- if self-hosters ever want a unique repo of users per site, they have to remove these constraints manually
    ,CONSTRAINT pm_user_username_key UNIQUE (username)
    ,CONSTRAINT pm_user_email_key UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS pm_user_authz (
    site_id UUID
    ,user_id UUID
    ,roles JSON
    ,authz_attributes JSON
    ,role_authz_attributes JSON

    ,CONSTRAINT pm_user_authz_site_id_user_id_pkey PRIMARY KEY (site_id, user_id)
    ,CONSTRAINT pm_user_authz_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id)
    ,CONSTRAINT pm_user_authz_user_id_fkey FOREIGN KEY (user_id) REFERENCES pm_user (user_id)
);

CREATE INDEX IF NOT EXISTS pm_user_authz_site_id_idx ON pm_user_authz (site_id);

CREATE INDEX IF NOT EXISTS pm_user_authz_user_id_idx ON pm_user_authz (user_id);

CREATE TRIGGER pm_user_authz_roles_after_insert AFTER INSERT ON pm_user_authz
BEGIN
    INSERT INTO pm_user_authz_roles_tblidx
        (site_id, user_id, role)
    SELECT
        NEW.site_id
        ,NEW.user_id
        ,json_each.value AS role
    FROM
        json_each(NEW.roles)
    WHERE
        TRUE
    ON CONFLICT DO NOTHING;
END;

CREATE TRIGGER pm_user_authz_roles_after_update AFTER UPDATE ON pm_user_authz
WHEN OLD.roles <> NEW.roles
BEGIN
    DELETE FROM pm_user_authz_roles_tblidx
    WHERE NOT EXISTS (
        SELECT 1
        FROM
            json_each(NEW.roles)
        WHERE
            (site_id, user_id) = (NEW.site_id, NEW.user_id)
            AND role = json_each.value
    );
    INSERT INTO pm_user_authz_roles_tblidx
        (site_id, user_id, role)
    SELECT
        NEW.site_id
        ,NEW.user_id
        ,json_each.value AS role
    FROM
        json_each(NEW.roles)
    WHERE
        TRUE
    ON CONFLICT DO NOTHING;
END;

CREATE TABLE IF NOT EXISTS pm_user_authz_roles_tblidx (
    site_id UUID
    ,user_id UUID
    ,role TEXT

    ,CONSTRAINT pm_user_authz_roles_tblidx_pkey PRIMARY KEY (site_id, user_id, role)
    ,CONSTRAINT pm_user_authz_roles_tblix_fkey FOREIGN KEY (site_id, user_id) REFERENCES pm_user_authz (site_id, user_id)
);

CREATE INDEX IF NOT EXISTS pm_user_authz_roles_tblix_site_id_user_id_idx ON pm_user_authz_roles_tblidx (site_id, user_id);

CREATE INDEX IF NOT EXISTS pm_user_authz_roles_tblidx_site_id_role_idx ON pm_user_authz_roles_tblidx (site_id, role);

CREATE TABLE IF NOT EXISTS pm_role (
    site_id UUID
    --,namespace TEXT
    ,role TEXT
    ,authz_attributes JSON

    ,CONSTRAINT pm_role_site_id_role_pkey PRIMARY KEY (site_id, role)
    ,CONSTRAINT pm_role_site_id_jkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id)
);

CREATE INDEX pm_role_site_id_idx ON pm_role (site_id);

CREATE TABLE pm_session (
    session_hash BLOB
    ,site_id UUID
    ,user_id UUID

    ,CONSTRAINT pm_session_session_hash_pkey PRIMARY KEY (session_hash)
    ,CONSTRAINT pm_session_site_id_user_id FOREIGN KEY (site_id, user_id) REFERENCES pm_user_authz (site_id, user_id)
    ,CONSTRAINT pm_session_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id)
    ,CONSTRAINT pm_session_user_id_fkey FOREIGN KEY (user_id) REFERENCES pm_user (user_id)
);

CREATE INDEX IF NOT EXISTS pm_session_site_id_user_id_idx ON pm_session (site_id, user_id);

CREATE INDEX IF NOT EXISTS pm_session_site_id_idx ON pm_session (site_id);

CREATE INDEX IF NOT EXISTS pm_session_user_id_idx ON pm_session (user_id);
