-- need administrate-sites
CREATE TABLE pm_site (
    site_id UUID NOT NULL
    ,domain TEXT NOT NULL
    ,subdomain TEXT NOT NULL
    ,tilde_prefix TEXT NOT NULL
    ,is_primary BOOLEAN DEFAULT NULL

    ,CONSTRAINT pm_site_site_id_pkey PRIMARY KEY (site_id)
    ,CONSTRAINT pm_site_domain_subdomain_tilde_prefix_key UNIQUE (domain, subdomain, tilde_prefix)
    ,CONSTRAINT pm_site_is_primary_key UNIQUE (is_primary)
);

-- handled on startup
CREATE TABLE pm_plugin (
    plugin TEXT NOT NULL
    ,url TEXT NOT NULL
    ,version TEXT NOT NULL

    ,CONSTRAINT pm_plugin_plugin_pkey PRIMARY KEY (plugin)
    ,CONSTRAINT pm_plugin_url_plugin_key UNIQUE (url, version)
);

-- handled on startup
CREATE TABLE pm_handler (
    plugin TEXT NOT NULL
    ,handler TEXT NOT NULL

    ,CONSTRAINT pm_handler_plugin_handler_pkey PRIMARY KEY (plugin, handler)
    ,CONSTRAINT pm_handler_plugin_fkey FOREIGN KEY (plugin) REFERENCES pm_plugin (plugin) ON UPDATE CASCADE
);

CREATE INDEX pm_handler_plugin_idx ON pm_handler (plugin);

-- handled on startup
CREATE TABLE pm_capability (
    plugin TEXT NOT NULL
    ,capability TEXT NOT NULL

    ,CONSTRAINT pm_capability_plugin_capability_pkey PRIMARY KEY (plugin, capability)
    ,CONSTRAINT pm_capability_plugin_fkey FOREIGN KEY (plugin) REFERENCES pm_plugin (plugin) ON UPDATE CASCADE
);

CREATE INDEX pm_capability_plugin_idx ON pm_capability (plugin);

-- need administrate-sites
CREATE TABLE pm_allowed_plugin (
    plugin TEXT NOT NULL
    ,site_id UUID NOT NULL

    ,CONSTRAINT pm_allowed_plugin_plugin_site_id_pkey PRIMARY KEY (plugin, site_id)
    ,CONSTRAINT pm_allowed_plugin_plugin_fkey FOREIGN KEY (plugin) REFERENCES pm_plugin (plugin) ON UPDATE CASCADE
    ,CONSTRAINT pm_allowed_plugin_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id) ON UPDATE CASCADE
);

CREATE INDEX pm_allowed_plugin_plugin_idx ON pm_allowed_plugin (plugin);

CREATE INDEX pm_allowed_plugin_site_id_idx ON pm_allowed_plugin (site_id);

-- need administrate-sites
CREATE TABLE pm_denied_plugin (
    plugin TEXT NOT NULL
    ,site_id UUID NOT NULL

    ,CONSTRAINT pm_denied_plugin_plugin_site_id_pkey PRIMARY KEY (plugin, site_id)
    ,CONSTRAINT pm_denied_plugin_plugin_fkey FOREIGN KEY (plugin) REFERENCES pm_plugin (plugin) ON UPDATE CASCADE
    ,CONSTRAINT pm_denied_plugin_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site) ON UPDATE CASCADE
);

CREATE INDEX pm_denied_plugin_plugin_idx ON pm_denied_plugin (plugin);

CREATE INDEX pm_denied_plugin_site_id_idx ON pm_denied_plugin (site_id);

-- need administrate-sites
CREATE TABLE pm_allowed_handler (
    plugin TEXT NOT NULL
    ,handler TEXT NOT NULL
    ,site_id UUID NOT NULL

    ,CONSTRAINT pm_allowed_handler_plugin_handler_site_id_pkey PRIMARY KEY (plugin, handler, site_id)
    ,CONSTRAINT pm_allowed_handler_plugin_handler_fkey FOREIGN KEY (plugin, handler) REFERENCES pm_handler (plugin, handler) ON UPDATE CASCADE
    ,CONSTRAINT pm_allowed_handler_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id) ON UPDATE CASCADE
);

CREATE INDEX pm_allowed_handler_plugin_handler_idx ON pm_allowed_handler (plugin, handler);

CREATE INDEX pm_allowed_handler_site_id_idx ON pm_allowed_handler (site_id);

-- need administrate-sites
CREATE TABLE pm_denied_handler (
    plugin TEXT NOT NULL
    ,handler TEXT NOT NULL
    ,site_id UUID NOT NULL

    ,CONSTRAINT pm_denied_handler_plugin_handler_site_id_pkey PRIMARY KEY (plugin, handler, site_id)
    ,CONSTRAINT pm_denied_handler_plugin_handler_fkey FOREIGN KEY (plugin, handler) REFERENCES pm_handler (plugin, handler) ON UPDATE CASCADE
    ,CONSTRAINT pm_denied_handler_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id) ON UPDATE CASCADE
);

CREATE INDEX pm_denied_handler_plugin_handler_idx ON pm_denied_handler (plugin, handler);

CREATE INDEX pm_denied_handler_site_id_idx ON pm_denied_handler (site_id);

-- need administrate-roles
CREATE TABLE pm_role (
    site_id UUID NOT NULL
    ,plugin TEXT NOT NULL
    ,role TEXT NOT NULL

    ,CONSTRAINT pm_role_site_id_plugin_role_pkey PRIMARY KEY (site_id, plugin, role)
    ,CONSTRAINT pm_role_site_id_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id) ON UPDATE CASCADE
    ,CONSTRAINT pm_role_plugin_fkey FOREIGN KEY (plugin) REFERENCES pm_plugin (plugin) ON UPDATE CASCADE
);

CREATE INDEX pm_role_site_id_idx ON pm_role (site_id);

CREATE INDEX pm_role_plugin_idx ON pm_role (plugin);

-- need administrate-tags
CREATE TABLE pm_tag (
    site_id UUID NOT NULL
    ,plugin TEXT NOT NULL
    ,tag TEXT NOT NULL

    ,CONSTRAINT pm_tag_site_id_plugin_tag_pkey PRIMARY KEY (site_id, plugin, tag)
    ,CONSTRAINT pm_tag_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id) ON UPDATE CASCADE
    ,CONSTRAINT pm_tag_plugin_fkey FOREIGN KEY (plugin) REFERENCES pm_plugin (plugin) ON UPDATE CASCADE
);

CREATE INDEX pm_tag_site_id_idx ON pm_tag (site_id);

CREATE INDEX pm_tag_plugin_idx ON pm_tag (plugin);

-- need administrate-roles
CREATE TABLE pm_role_capability (
    site_id UUID NOT NULL
    ,plugin TEXT NOT NULL
    ,role TEXT NOT NULL
    ,capability TEXT NOT NULL

    ,CONSTRAINT pm_role_capability_site_id_plugin_role_capability_pkey PRIMARY KEY (site_id, plugin, role, capability)
    ,CONSTRAINT pm_role_capability_site_id_plugin_role_fkey FOREIGN KEY (site_id, plugin, role) REFERENCES pm_role (site_id, plugin, role) ON UPDATE CASCADE
    ,CONSTRAINT pm_role_capability_plugin_capability_fkey FOREIGN KEY (plugin, capability) REFERENCES pm_capability (plugin, capability) ON UPDATE CASCADE
);

CREATE INDEX pm_role_capability_site_id_plugin_role_idx ON pm_role_capability (site_id, plugin, role);

CREATE INDEX pm_role_capability_plugin_capability_idx ON pm_role_capability (plugin, capability);

-- need administrate-tags
CREATE TABLE pm_tag_capability (
    site_id UUID NOT NULL
    ,plugin TEXT NOT NULL
    ,tag TEXT NOT NULL
    ,role TEXT NOT NULL
    ,capability TEXT NOT NULL

    ,CONSTRAINT pm_tag_capability_site_id_plugin_tag_role_pkey PRIMARY KEY (site_id, plugin, tag, role)
    ,CONSTRAINT pm_tag_capability_site_id_plugin_tag_fkey FOREIGN KEY (site_id, plugin, tag) REFERENCES pm_tag (site_id, plugin, tag) ON UPDATE CASCADE
    ,CONSTRAINT pm_tag_capability_site_id_plugin_role_fkey FOREIGN KEY (site_id, plugin, role) REFERENCES pm_role (site_id, plugin, role) ON UPDATE CASCADE
    ,CONSTRAINT pm_tag_capability_plugin_capability_fkey FOREIGN KEY (plugin, capability) REFERENCES pm_capability (plugin, capability) ON UPDATE CASCADE
);

CREATE INDEX pm_tag_capability_site_id_plugin_tag_idx ON pm_tag_capability (site_id, plugin, tag);

CREATE INDEX pm_tag_capability_site_id_plugin_role_idx ON pm_tag_capability (site_id, plugin, role);

CREATE INDEX pm_tag_capability_plugin_capability_idx ON pm_tag_capability (plugin, capability);

-- need administrate-tags
CREATE TABLE pm_tag_owner (
    site_id UUID NOT NULL
    ,plugin TEXT NOT NULL
    ,tag TEXT NOT NULL
    ,role TEXT NOT NULL

    ,CONSTRAINT pm_tag_owner_site_id_plugin_tag_role_pkey PRIMARY KEY (site_id, plugin, tag, role)
    ,CONSTRAINT pm_tag_owner_site_id_plugin_tag_fkey FOREIGN KEY (site_id, plugin, tag) REFERENCES pm_tag (site_id, plugin, tag) ON UPDATE CASCADE
    ,CONSTRAINT pm_tag_owner_site_id_plugin_role_fkey FOREIGN KEY (site_id, plugin, role) REFERENCES pm_role (site_id, plugin, role) ON UPDATE CASCADE
);

CREATE INDEX pm_tag_owner_site_id_plugin_tag_idx ON pm_tag_owner (site_id, plugin, tag);

CREATE INDEX pm_tag_owner_site_id_plugin_role_idx ON pm_tag_owner (site_id, plugin, role);

CREATE TABLE pm_user (
    site_id UUID NOT NULL
    ,user_id UUID NOT NULL
    ,name TEXT
    ,username TEXT NOT NULL
    ,email TEXT NOT NULL
    ,password_hash TEXT
    ,reset_password_token TEXT
    ,reset_password_sent_at DATETIME

    ,CONSTRAINT pm_user_site_id_user_id_pkey PRIMARY KEY (site_id, user_id)
    ,CONSTRAINT pm_user_site_id_username_key UNIQUE (site_id, username)
    ,CONSTRAINT pm_user_site_id_email_key UNIQUE (site_id, email)
    ,CONSTRAINT pm_user_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id) ON UPDATE CASCADE
);

CREATE INDEX pm_user_site_id_idx ON pm_user (site_id);

CREATE TABLE pm_session (
    site_id UUID NOT NULL
    ,session_hash BLOB NOT NULL
    ,user_id UUID NOT NULL

    ,CONSTRAINT pm_session_site_id_session_hash_pkey PRIMARY KEY (site_id, session_hash)
    ,CONSTRAINT pm_session_site_id_user_id_fkey FOREIGN KEY (site_id, user_id) REFERENCES pm_user (site_id, user_id) ON UPDATE CASCADE
);

CREATE INDEX pm_session_site_id_user_id_idx ON pm_session (site_id, user_id);

-- need assign-roles
CREATE TABLE pm_user_role (
    site_id UUID NOT NULL
    ,user_id UUID NOT NULL
    ,plugin TEXT NOT NULL
    ,role TEXT NOT NULL

    ,CONSTRAINT pm_user_role_site_id_user_id_plugin_role_pkey PRIMARY KEY (site_id, user_id, plugin, role)
    ,CONSTRAINT pm_user_role_site_id_user_id_fkey FOREIGN KEY (site_id, user_id) REFERENCES pm_user (site_id, user_id) ON UPDATE CASCADE
    ,CONSTRAINT pm_user_role_site_id_plugin_role_fkey FOREIGN KEY (site_id, plugin, role) REFERENCES pm_role (site_id, plugin, role) ON UPDATE CASCADE
);

CREATE INDEX pm_user_role_site_id_user_id_fkey ON pm_user_role (site_id, user_id);

CREATE INDEX pm_user_role_site_id_plugin_role_idx ON pm_user_role (site_id, plugin, role);

-- need url.edit-all or url.edit-all-handlers or url.edit-all-handler-configs
-- for the corresponding URL
CREATE TABLE pm_url (
    site_id UUID NOT NULL
    ,urlpath TEXT NOT NULL
    ,plugin TEXT NOT NULL
    ,handler TEXT NOT NULL
    ,config JSON

    ,CONSTRAINT pm_url_site_id_urlpath_pkey PRIMARY KEY (site_id, urlpath)
    ,CONSTRAINT pm_url_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id) ON UPDATE CASCADE
    ,CONSTRAINT pm_url_plugin_handler_fkey FOREIGN KEY (plugin, handler) REFERENCES pm_handler (plugin, handler) ON UPDATE CASCADE
);

CREATE INDEX pm_url_site_id_idx ON pm_url (site_id);

CREATE INDEX pm_url_plugin_handler_idx ON pm_url (plugin, handler);

-- need url.administrate capability for the corresponding URL
CREATE TABLE pm_url_role_capability (
    site_id UUID NOT NULL
    ,urlpath TEXT NOT NULL
    ,plugin TEXT NOT NULL
    ,role TEXT NOT NULL
    ,capability TEXT NOT NULL

    ,CONSTRAINT pm_url_role_capability_site_id_urlpath_plugin_role_capabil_pkey PRIMARY KEY (site_id, urlpath, plugin, role, capability)
    ,CONSTRAINT pm_url_role_capability_site_id_plugin_role_fkey FOREIGN KEY (site_id, plugin, role) REFERENCES pm_role (site_id, plugin, role) ON UPDATE CASCADE
    ,CONSTRAINT pm_url_role_capability_plugin_capability_fkey FOREIGN KEY (plugin, capability) REFERENCES pm_capability (plugin, capability) ON UPDATE CASCADE
);

CREATE INDEX pm_url_role_capability_site_id_plugin_role_idx ON pm_url_role_capability (site_id, plugin, role);

CREATE INDEX pm_url_role_capability_plugin_capability_idx ON pm_url_role_capability (plugin, capability);

-- need url.administrate capability for the corresponding URL
CREATE TABLE pm_url_tag (
    site_id UUID NOT NULL
    ,urlpath TEXT NOT NULL
    ,plugin TEXT NOT NULL
    ,tag TEXT NOT NULL

    ,CONSTRAINT pm_url_tag_site_id_urlpath_plugin_tag_pkey PRIMARY KEY (site_id, urlpath, plugin, tag)
    ,CONSTRAINT pm_url_tag_site_id_plugin_tag_fkey FOREIGN KEY (site_id, plugin, tag) REFERENCES pm_tag (site_id, plugin, tag) ON UPDATE CASCADE
);

CREATE INDEX pm_url_tag_site_id_plugin_tag_idx ON pm_url_tag (site_id, plugin, tag);

-- handled by pagemanager/template
CREATE TABLE pm_template_data (
    site_id UUID NOT NULL
    ,langcode TEXT NOT NULL
    ,data_file TEXT NOT NULL
    ,data JSON

    ,CONSTRAINT pm_template_data_site_id_langcode_data_file_pkey PRIMARY KEY (site_id, langcode, data_file)
    ,CONSTRAINT pm_template_data_site_id_fkey FOREIGN KEY (site_id) REFERENCES pm_site (site_id) ON UPDATE CASCADE
);

CREATE INDEX pm_template_data_site_id_idx ON pm_template_data (site_id);
