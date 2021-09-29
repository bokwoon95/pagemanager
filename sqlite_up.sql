CREATE TABLE IF NOT EXISTS pm_site (
    site_id UUID
    ,domain TEXT
    ,subdomain TEXT

    ,CONSTRAINT pm_site_site_id_pkey PRIMARY KEY (site_id)
    ,CONSTRAINT pm_site_domain_subdomain_key UNIQUE (domain, subdomain)
);

-- url -> domain, subdomain, langcode, path -> site_id, langcode, path
CREATE TABLE IF NOT EXISTS pm_url (
    site_id UUID
    ,path TEXT
    ,plugin TEXT
    ,handler TEXT
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

CREATE TABLE IF NOT EXISTS pm_role (
    site_id UUID
    ,role TEXT

    ,CONSTRAINT pm_role_site_id_role PRIMARY KEY (site_id, role)
    ,CONSTRAINT pm_role_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id)
);

CREATE INDEX IF NOT EXISTS pm_role_site_id_idx ON pm_role (site_id);

CREATE TABLE IF NOT EXISTS pm_user_role (
    site_id UUID
    ,user_id UUID
    ,role TEXT

    ,CONSTRAINT pm_user_role_site_id_user_id_role_pkey PRIMARY KEY (site_id, user_id, role)
    ,CONSTRAINT pm_user_role_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id)
    ,CONSTRAINT pm_user_role_user_id_fkey FOREIGN KEY (user_id) REFERENCES pm_user (user_id)
);

CREATE INDEX IF NOT EXISTS pm_user_role_site_id ON pm_user_role (site_id);

CREATE INDEX IF NOT EXISTS pm_user_role_user_id ON pm_user_role (user_id);

CREATE TABLE IF NOT EXISTS pm_policy (
    site_id UUID
    ,role TEXT
    ,label TEXT -- pm-superadmin, pm-admin, pm-url
    ,action TEXT
);

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
