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

CREATE TABLE IF NOT EXISTS pm_role_user (
    site_id UUID
    ,role TEXT
    ,user_id UUID

    ,CONSTRAINT pm_role_user_site_id_role_user_id_pkey PRIMARY KEY (site_id, role, user_id)
    ,CONSTRAINT pm_role_user_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id)
    ,CONSTRAINT pm_role_user_user_id_fkey FOREIGN KEY (user_id) REFERENCES pm_user (user_id)
);

CREATE INDEX IF NOT EXISTS pm_role_user_site_id ON pm_role_user (site_id);

CREATE INDEX IF NOT EXISTS pm_role_user_user_id ON pm_role_user (user_id);

-- site_id, label

CREATE TABLE IF NOT EXISTS pm_permission (
    site_id UUID
    ,role TEXT
    ,label TEXT -- pm-superadmin, pm-admin, pm-url
    ,action TEXT

    ,CONSTRAINT pm_permission_site_id_role_label_action_pkey PRIMARY KEY (site_id, role, label, action)
    ,CONSTRAINT pm_permission_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id)
    ,CONSTRAINT pm_permission_site_id_role_fkey FOREIGN KEY (site_id, role) REFERENCES pm_role (site_id, role)
);

CREATE INDEX IF NOT EXISTS pm_permission_site_id_idx ON pm_permission (site_id);

CREATE INDEX IF NOT EXISTS pm_permission_site_id_role_idx ON pm_permission (site_id, role);

CREATE TABLE IF NOT EXISTS pm_session (
    session_hash BLOB
    ,site_id UUID
    ,user_id UUID

    ,CONSTRAINT pm_session_session_hash_pkey PRIMARY KEY (session_hash)
    ,CONSTRAINT pm_session_site_id_user_id_fkey FOREIGN KEY (site_id, user_id) REFERENCES pm_user_authz (site_id, user_id)
    ,CONSTRAINT pm_session_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id)
    ,CONSTRAINT pm_session_user_id_fkey FOREIGN KEY (user_id) REFERENCES pm_user (user_id)
);

CREATE INDEX IF NOT EXISTS pm_session_site_id_user_id_idx ON pm_session (site_id, user_id);

CREATE INDEX IF NOT EXISTS pm_session_site_id_idx ON pm_session (site_id);

CREATE INDEX IF NOT EXISTS pm_session_user_id_idx ON pm_session (user_id);
