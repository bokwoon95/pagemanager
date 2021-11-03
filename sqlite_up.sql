-- need administrate-sites
CREATE TABLE IF NOT EXISTS pm_site (
    site_id UUID
    ,domain TEXT NOT NULL
    ,subdomain TEXT NOT NULL
    ,tilde_prefix TEXT NOT NULL
    ,is_primary BOOLEAN DEFAULT NULL

    ,CONSTRAINT pm_site_site_id_pkey PRIMARY KEY (site_id)
    ,CONSTRAINT pm_site_domain_subdomain_tilde_prefix_key UNIQUE (domain, subdomain, tilde_prefix)
    ,CONSTRAINT pm_site_is_primary_key UNIQUE (is_primary)
);

-- need administrate-sites
CREATE TABLE IF NOT EXISTS pm_site_plugin (
    plugin TEXT
    ,allowed_sites JSON
    ,denied_sites JSON

    ,CONSTRAINT pm_site_plugin_pkey PRIMARY KEY (plugin)
);

-- need administrate-sites
CREATE TABLE IF NOT EXISTS pm_site_handler (
    plugin TEXT
    ,handler TEXT
    ,allowed_sites JSON
    ,denied_sites JSON

    ,CONSTRAINT pm_site_handler_plugin_handler_pkey PRIMARY KEY (plugin, handler)
);

CREATE TABLE IF NOT EXISTS pm_user (
    site_id UUID
    ,user_id UUID
    ,name TEXT
    ,username TEXT NOT NULL
    ,email TEXT NOT NULL
    ,password_hash TEXT
    ,reset_password_token TEXT
    ,reset_password_sent_at DATETIME

    ,CONSTRAINT pm_user_site_id_user_id_pkey PRIMARY KEY (site_id, user_id)
    ,CONSTRAINT pm_user_site_id_username_key UNIQUE (site_id, username)
    ,CONSTRAINT pm_user_site_id_email_key UNIQUE (site_id, email)
);

CREATE TABLE IF NOT EXISTS pm_session (
    site_id UUID
    ,session_hash BLOB
    ,user_id UUID NOT NULL

    ,CONSTRAINT pm_session_site_id_session_hash_pkey PRIMARY KEY (site_id, session_hash)
    ,CONSTRAINT pm_session_site_id_user_id_fkey FOREIGN KEY (site_id, user_id) REFERENCES pm_user (site_id, user_id)
);

CREATE INDEX pm_session_site_id_user_id_idx ON pm_session (site_id, user_id);

-- need administrate-roles
CREATE TABLE IF NOT EXISTS pm_role (
    site_id UUID
    ,plugin TEXT
    ,role TEXT

    ,CONSTRAINT pm_role_site_id_plugin_role_pkey PRIMARY KEY (site_id, plugin, role)
);

-- need administrate-tags
CREATE TABLE IF NOT EXISTS pm_tag (
    site_id UUID
    ,plugin TEXT
    ,tag TEXT

    ,CONSTRAINT pm_tag_site_id_plugin_tag_pkey PRIMARY KEY (site_id, plugin, tag)
);

-- need administrate-roles
CREATE TABLE IF NOT EXISTS pm_role_capability (
    site_id UUID
    ,plugin TEXT
    ,role TEXT
    ,capability TEXT

    ,CONSTRAINT pm_role_capability_site_id_plugin_role_capability_pkey PRIMARY KEY (site_id, plugin, role, capability)
    ,CONSTRAINT pm_role_capability_site_id_plugin_role_fkey FOREIGN KEY (site_id, plugin, role) REFERENCES pm_role (site_id, plugin, role)
);

CREATE INDEX pm_role_capability_site_id_plugin_role_idx ON pm_role_capability (site_id, plugin, role);

-- need assign-roles
CREATE TABLE IF NOT EXISTS pm_user_role (
    site_id UUID
    ,plugin TEXT
    ,user_id UUID
    ,role TEXT

    ,CONSTRAINT pm_user_role_site_id_plugin_user_id_role_pkey PRIMARY KEY (site_id, plugin, user_id, role)
    ,CONSTRAINT pm_user_role_site_id_plugin_role_fkey FOREIGN KEY (site_id, plugin, role) REFERENCES pm_role (site_id, plugin, role)
);

CREATE INDEX pm_user_role_site_id_plugin_role_idx ON pm_user_role (site_id, plugin, role);

-- need administrate-tags
CREATE TABLE IF NOT EXISTS pm_tag_capability (
    site_id UUID
    ,plugin TEXT
    ,tag TEXT
    ,role TEXT
    ,capability TEXT NOT NULL

    ,CONSTRAINT pm_tag_capability_site_id_plugin_tag_role_pkey PRIMARY KEY (site_id, plugin, tag, role)
    ,CONSTRAINT pm_tag_capability_site_id_plugin_tag_fkey FOREIGN KEY (site_id, plugin, tag) REFERENCES pm_tag (site_id, plugin, tag)
);

CREATE INDEX pm_tag_capability_site_id_plugin_tag_idx ON pm_tag_capability (site_id, plugin, tag);

-- need administrate-tags
CREATE TABLE IF NOT EXISTS pm_tag_owner (
    site_id UUID
    ,plugin TEXT
    ,tag TEXT
    ,role TEXT

    ,CONSTRAINT pm_tag_owner_site_id_plugin_tag_role_pkey PRIMARY KEY (site_id, plugin, tag, role)
    ,CONSTRAINT pm_tag_owner_site_id_plugin_tag_fkey FOREIGN KEY (site_id, plugin, tag) REFERENCES pm_tag (site_id, plugin, tag)
    ,CONSTRAINT pm_tag_owner_site_id_plugin_role_fkey FOREIGN KEY (site_id, plugin, role) REFERENCES pm_role (site_id, plugin, role)
);

CREATE INDEX pm_tag_owner_site_id_plugin_tag_idx ON pm_tag_owner (site_id, plugin, tag);

CREATE INDEX pm_tag_owner_site_id_plugin_role_idx ON pm_tag_owner (site_id, plugin, role);

-- need url.edit-all or url.edit-all-handlers or url.edit-all-handler-configs
-- for the corresponding URL
CREATE TABLE IF NOT EXISTS pm_url (
    site_id UUID
    ,urlpath TEXT
    ,plugin TEXT NOT NULL
    ,handler TEXT NOT NULL
    ,config JSON

    ,CONSTRAINT pm_url_site_id_urlpath_pkey PRIMARY KEY (site_id, urlpath)
);

CREATE INDEX pm_url_site_id_idx ON pm_url (site_id);

-- need url.administrate capability for the corresponding URL
-- plugin is implicitly ''
CREATE TABLE IF NOT EXISTS pm_url_role_capability (
    site_id UUID
    ,urlpath TEXT
    ,role TEXT
    ,capability TEXT

    ,CONSTRAINT pm_url_role_capability_site_id_urlpath_role_capability_pkey PRIMARY KEY (site_id, urlpath, role, capability)
);

-- need url.administrate capability for the corresponding URL
-- plugin is implicitly ''
CREATE TABLE IF NOT EXISTS pm_url_tag (
    site_id UUID
    ,urlpath TEXT
    ,tag TEXT

    ,CONSTRAINT pm_url_tag_site_id_urlpath_tag_pkey PRIMARY KEY (site_id, urlpath, tag)
);

-- handled by pagemanager/blog
CREATE TABLE IF NOT EXISTS pm_template_data (
    site_id UUID
    ,langcode TEXT
    ,data_file TEXT
    ,data JSON

    ,CONSTRAINT pm_template_data_site_id_langcode_data_file_pkey PRIMARY KEY (site_id, langcode, data_file)
);
