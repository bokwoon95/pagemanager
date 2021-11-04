INSERT INTO pm_site
    (site_id, domain, subdomain, tilde_prefix, is_primary)
VALUES
    (x'11111111111111111111111111111111', 'localhost', '', '', TRUE)
ON CONFLICT DO NOTHING
;

INSERT INTO pm_plugin
    (plugin, url, version)
VALUES
    ('pagemanager', 'github.com/pagemanager/pagemanager', '')
    ,('template', 'github.com/pagemanager/template', '')
ON CONFLICT DO NOTHING
;

INSERT INTO pm_handler
    (plugin, handler)
VALUES
    ('pagemanager', 'url-dashboard')
ON CONFLICT DO NOTHING
;

INSERT INTO pm_user
    (site_id, user_id, username, email)
VALUES
    (x'11111111111111111111111111111111', x'11111111111111111111111111111111', 'superadmin', 'superadmin@email.com')
ON CONFLICT DO NOTHING
;

INSERT INTO pm_capability
    (plugin, capability)
VALUES
    ('pagemanager', 'administrate-sites')
    ,('pagemanager', 'administrate-roles')
    ,('pagemanager', 'administrate-tags')
    ,('pagemanager', 'assign-roles')
    ,('pagemanager', 'url.administrate')
    ,('pagemanager', 'url.edit-all-entries')
    ,('pagemanager', 'url.edit-all-handlers')
    ,('pagemanager', 'url.edit-all-handler-configs')
    ,('pagemanager', 'url.edit-handler')
    ,('pagemanager', 'url.edit-handler-config')
ON CONFLICT DO NOTHING
;

INSERT INTO pm_role
    (site_id, plugin, role)
VALUES
    (x'11111111111111111111111111111111', 'pagemanager', 'superadmin')
    ,(x'11111111111111111111111111111111', 'pagemanager', 'admin')
ON CONFLICT DO NOTHING
;

INSERT INTO pm_role_capability
    (site_id, plugin, role, capability)
VALUES
    (x'11111111111111111111111111111111', 'pagemanager', 'superadmin', 'administrate-sites')
    ,(x'11111111111111111111111111111111', 'pagemanager', 'superadmin', 'administrate-roles')
    ,(x'11111111111111111111111111111111', 'pagemanager', 'superadmin', 'administrate-tags')
    ,(x'11111111111111111111111111111111', 'pagemanager', 'admin', 'assign-roles')
ON CONFLICT DO NOTHING
;

INSERT INTO pm_user
    (site_id, user_id, username, email)
VALUES
    (x'11111111111111111111111111111111', x'11111111111111111111111111111111', 'admin', 'admin@localhost')
ON CONFLICT DO NOTHING
;

INSERT INTO pm_user_role
    (site_id, user_id, plugin, role)
VALUES
    (x'11111111111111111111111111111111', x'11111111111111111111111111111111', 'pagemanager', 'superadmin')
    ,(x'11111111111111111111111111111111', x'11111111111111111111111111111111', 'pagemanager', 'admin')
ON CONFLICT DO NOTHING
;

INSERT INTO pm_url
    (site_id, urlpath, plugin, handler, config)
VALUES
    (x'11111111111111111111111111111111', '/pm-url', 'pagemanager', 'url-dashboard', NULL)
ON CONFLICT DO NOTHING
;

INSERT INTO pm_url_role_capability
    (site_id, urlpath, plugin, role, capability)
VALUES
    (x'11111111111111111111111111111111', '/', 'pagemanager', 'admin', 'url.administrate')
    ,(x'11111111111111111111111111111111', '/', 'pagemanager', 'admin', 'url.edit-all-entries')
ON CONFLICT DO NOTHING
;

-- I want to know what capabilities a user has for a url
-- get all the capabilities on url '/' for site_id = x'11111111111111111111111111111111' and user_id = x'11111111111111111111111111111111'
SELECT
    pm_url_role_capability.capability
FROM
    pm_url_role_capability
    JOIN pm_user_role USING (site_id, plugin, role)
    JOIN pm_plugin USING (plugin)
WHERE
    pm_url_role_capability.site_id = x'11111111111111111111111111111111'
    AND pm_user_role.user_id = x'11111111111111111111111111111111'
    AND pm_plugin.url = 'github.com/pagemanager/pagemanager'
    AND pm_plugin.version = ''
    AND pm_url_role_capability.urlpath = '/'
;

-- I want to know what urls a user can edit
-- get all the urls for site_id = x'11111111111111111111111111111111' where user_id = x'11111111111111111111111111111111' has capability 'url.edit-all-entries'
SELECT
    pm_url_role_capability.urlpath
FROM
    pm_url_role_capability
    JOIN pm_user_role USING (site_id, plugin, role)
    JOIN pm_plugin USING (plugin)
WHERE
    pm_url_role_capability.site_id = x'11111111111111111111111111111111'
    AND pm_user_role.user_id = x'11111111111111111111111111111111'
    AND pm_plugin.url = 'github.com/pagemanager/pagemanager'
    AND pm_plugin.version = ''
    AND pm_url_role_capability.capability = 'url.edit-all-entries'
;
