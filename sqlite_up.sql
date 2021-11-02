CREATE TABLE IF NOT EXISTS pm_site (
    site_id UUID
    ,domain TEXT NOT NULL
    ,subdomain TEXT NOT NULL
    ,tilde_prefix TEXT NOT NULL
    ,created_at DATETIME DEFAULT CURRENT_TIMESTAMP -- earliest site is default site used in 'offline'
    ,updated_at DATETIME DEFAULT CURRENT_TIMESTAMP

    ,CONSTRAINT pm_site_site_id_pkey PRIMARY KEY (site_id)
    ,CONSTRAINT pm_site_domain_subdomain_tilde_prefix_key UNIQUE (domain, subdomain, tilde_prefix)
);

CREATE TABLE IF NOT EXISTS pm_user (
    user_id UUID
    ,email TEXT
    ,username TEXT
    ,password_hash TEXT
    ,name TEXT
    ,created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    ,updated_at DATETIME DEFAULT CURRENT_TIMESTAMP

    ,CONSTRAINT pm_user_user_id_pkey PRIMARY KEY (user_id)
    ,CONSTRAINT pm_user_email_key UNIQUE (email)
);

-- up to the plugin whether they want to use this table
CREATE TABLE IF NOT EXISTS pm_site_user (
    site_id UUID
    ,user_id UUID

    ,CONSTRAINT pm_site_user_pkey PRIMARY KEY (site_id, user_id)
);

-- need assign_roles
CREATE TABLE IF NOT EXISTS pm_user_role (
    site_id UUID
    ,plugin TEXT
    ,user_id UUID
    ,role TEXT

    ,CONSTRAINT pm_user_role_site_id_plugin_user_id_role_pkey PRIMARY KEY (site_id, plugin, user_id, role)
);

-- need administrate_roles
CREATE TABLE IF NOT EXISTS pm_role_capability (
    site_id UUID
    ,plugin TEXT
    ,role TEXT
    ,capability TEXT

    ,CONSTRAINT pm_role_capability_site_id_plugin_role_capability PRIMARY KEY (site_id, plugin, role, capability)
);

-- need administrate_tags
CREATE TABLE IF NOT EXISTS pm_tag_capability (
    site_id UUID
    ,plugin TEXT
    ,tag TEXT
    ,role TEXT
    ,capability TEXT NOT NULL

    ,CONSTRAINT pm_tag_capability_site_id_plugin_tag_role_pkey PRIMARY KEY (site_id, plugin, tag, role)
);

-- need administrate_tags
CREATE TABLE IF NOT EXISTS pm_tag_owner (
    site_id UUID
    ,plugin TEXT
    ,tag TEXT
    ,role TEXT

    ,CONSTRAINT pm_tag_owner_site_id_plugin_tag_role_pkey PRIMARY KEY (site_id, plugin, tag, role)
);

-- need edit_url_entries/edit_handlers/edit_handler_configs for that URL
CREATE TABLE IF NOT EXISTS pm_url (
    site_id UUID
    ,urlpath TEXT
    ,plugin TEXT NOT NULL
    ,handler TEXT NOT NULL
    ,config JSON

    ,CONSTRAINT pm_url_site_id_urlpath_pkey PRIMARY KEY (site_id, urlpath)
);

CREATE INDEX pm_url_site_id_idx ON pm_url (site_id);

-- need administrate_url capability for that URL
CREATE TABLE IF NOT EXISTS pm_url_role_capability (
    site_id UUID
    ,urlpath TEXT
    ,role TEXT
    ,capability TEXT

    ,CONSTRAINT pm_url_role_capability_site_id_urlpath_role_capability_pkey PRIMARY KEY (site_id, urlpath, role, capability)
);

-- need administrate_url capability for that URL
CREATE TABLE IF NOT EXISTS pm_url_tag (
    site_id UUID
    ,urlpath TEXT
    ,tag TEXT

    ,CONSTRAINT pm_url_tag_site_id_urlpath_tag_pkey PRIMARY KEY (site_id, urlpath, tag)
);

-- handled by plugin
CREATE TABLE IF NOT EXISTS pm_template_data (
    site_id UUID
    ,langcode TEXT
    ,data_file TEXT
    ,data JSON

    ,CONSTRAINT pm_template_data_site_id_langcode_data_file_pkey PRIMARY KEY (site_id, langcode, data_file)
);

-- handled by application
CREATE TABLE IF NOT EXISTS pm_session (
    session_hash BLOB
    ,user_id UUID NOT NULL

    ,CONSTRAINT pm_session_session_hash_pkey PRIMARY KEY (session_hash)
);
