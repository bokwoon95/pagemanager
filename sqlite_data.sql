INSERT INTO pm_site
    (site_id)
VALUES
    (0)
ON CONFLICT DO NOTHING
;

INSERT INTO pm_user
    (user_id)
VALUES
    (0)
ON CONFLICT DO NOTHING
;

INSERT INTO pm_user_authz
    (site_id, user_id, roles)
VALUES
    (0, 0, '["alpha","beta","delta"]')
ON CONFLICT DO NOTHING
;
