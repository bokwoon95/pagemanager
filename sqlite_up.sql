CREATE TABLE IF NOT EXISTS pm_site (
    site_id UUID
    ,domain TEXT NOT NULL
    ,subdomain TEXT NOT NULL
    ,path_prefix TEXT NOT NULL

    ,CONSTRAINT pm_site_site_id_pkey PRIMARY KEY (site_id)
    ,CONSTRAINT pm_site_domain_subdomain_path_prefix_key UNIQUE (domain, subdomain, path_prefix)
);

-- url (decomposes to ->) domain, subdomain, path_prefix, langcode, url_path (translates to ->) site_id, langcode, url_path
CREATE TABLE IF NOT EXISTS pm_url (
    site_id UUID
    ,url_path TEXT
    ,plugin TEXT
    ,handler TEXT
    ,params JSON

    ,CONSTRAINT pm_url_site_id_url_path_pkey PRIMARY KEY (site_id, url_path)
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

-- /user/$username-$base32_user_id
-- only $base32_user_id is used to lookup. If the $username doesn't match, initiate a redirect.

CREATE TABLE IF NOT EXISTS pm_user (
    user_id UUID
    ,email TEXT
    ,username TEXT
    ,name TEXT
    ,password_hash TEXT

    ,CONSTRAINT pm_user_user_id_pkey PRIMARY KEY (user_id)
    ,CONSTRAINT pm_user_email_key UNIQUE (email)
);

-- username is not unique by default: it is up to plugins to decide whether they want to enforce username uniqueness on a global level or on a site level.

-- pm_site_user is only used if the plugin requires maintaining a separate list of users per site. A plugin can always treat the entirety of pm_user as the site's install base.
CREATE TABLE IF NOT EXISTS pm_site_user (
    site_id UUID
    ,user_id UUID

    ,CONSTRAINT pm_site_user_site_id_user_id_pkey PRIMARY KEY (site_id, user_id)
    ,CONSTRAINT pm_site_user_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id)
    ,CONSTRAINT pm_site_user_user_id_fkey FOREIGN KEY (user_id) REFERENCES pm_user (user_id)
);

CREATE INDEX IF NOT EXISTS pm_site_user_site_id_idx ON pm_site_user (site_id);

CREATE INDEX IF NOT EXISTS pm_site_user_user_id_idx ON pm_site_user (user_id);

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

-- 0000 CRUD
-- 1111 15
-- 0100 4

-- 1 (1) CREATE
-- 10 (2) READ
-- 100 (4) UPDATE
-- 1000 (8) DELETE

-- READ WRITE DELETE

-- any object has an owner (user_id UUID) and object tags
-- admins have CRUD to all objects within the website

-- namespace:role
-- namespace:object_tag
-- namespace:object_tag:object_id

-- there's a table that maps object types to object tags
-- there's a table that maps objects to object types

-- roles users object_types objects
-- superadmin, admin

CREATE TABLE IF NOT EXISTS pm_permission (
    site_id UUID
    ,role TEXT
    ,label TEXT -- e.g. pm_url (should this be its own separate table? pm_resource?)
    --, label_params BLOB ?
    ,operation INT

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
    ,CONSTRAINT pm_session_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id)
    ,CONSTRAINT pm_session_user_id_fkey FOREIGN KEY (user_id) REFERENCES pm_user (user_id)
);

CREATE INDEX IF NOT EXISTS pm_session_site_id_idx ON pm_session (site_id);

CREATE INDEX IF NOT EXISTS pm_session_user_id_idx ON pm_session (user_id);
