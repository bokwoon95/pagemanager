INSERT INTO pm_site
    (site_id, domain, subdomain, tilde_prefix, is_primary)
VALUES
    (x'11111111111111111111111111111111', '', '', '', TRUE)
ON CONFLICT DO NOTHING
;

INSERT INTO pm_user
    (user_id)
VALUES
    (x'11111111111111111111111111111111')
ON CONFLICT DO NOTHING
;

INSERT INTO pm_user_role
    (site_id, plugin, user_id, role)
VALUES
    (x'11111111111111111111111111111111', '', x'11111111111111111111111111111111', 'superadmin')
    ,(x'11111111111111111111111111111111', '', x'11111111111111111111111111111111', 'admin')
ON CONFLICT DO NOTHING
;

INSERT INTO pm_role_capability
    (site_id, plugin, role, capability)
VALUES
    (x'11111111111111111111111111111111', '', 'superadmin', 'administrate_tags')
    ,(x'11111111111111111111111111111111', '', 'superadmin', 'administrate_roles')
    ,(x'11111111111111111111111111111111', '', 'admin', 'assign_roles')
ON CONFLICT DO NOTHING
;

INSERT INTO pm_url
    (site_id, urlpath, plugin, handler, config)
VALUES
    (x'11111111111111111111111111111111', '/pm-url', '', 'url_dashboard', NULL)
ON CONFLICT DO NOTHING
;

INSERT INTO pm_url_role_capability
    (site_id, urlpath, role, capability)
VALUES
    (x'11111111111111111111111111111111', '/', 'admin', 'administrate_url')
    ,(x'11111111111111111111111111111111', '/', 'admin', 'edit_url_entries')
ON CONFLICT DO NOTHING
;

-- get all the capabilities on url '/' for site_id = x'11111111111111111111111111111111' and user_id = x'11111111111111111111111111111111'
SELECT
    capability
FROM
    pm_url_role_capability AS a
    JOIN pm_user_role AS b USING (site_id, role)
WHERE
    a.site_id = x'11111111111111111111111111111111'
    AND b.user_id = x'11111111111111111111111111111111'
    AND urlpath = '/'
;

-- get all the urls for site_id = x'11111111111111111111111111111111' where user_id = x'11111111111111111111111111111111' has capability 'edit_url_entries'
