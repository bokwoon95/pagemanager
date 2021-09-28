package tables

import (
	"github.com/bokwoon95/pagemanager/sq"
)

type PM_SITE struct {
	sq.TableInfo `ddl:"unique={. cols=domain,subdomain}"`
	SITE_ID      sq.UUIDField `ddl:"primarykey"`
	DOMAIN       sq.StringField
	SUBDOMAIN    sq.StringField
}

type PM_URL struct {
	sq.TableInfo `ddl:"primarykey={. cols=site_id,url_path}"`
	SITE_ID      sq.UUIDField
	PATH         sq.StringField
	PLUGIN       sq.StringField
	PARAMS       sq.JSONField
}

type PM_TEMPLATE_DATA struct {
	sq.TableInfo `ddl:"primarykey={. cols=site_id,langcode,data_file}"`
	SITE_ID      sq.UUIDField
	LANGCODE     sq.StringField
	DATA_FILE    sq.StringField
	DATA         sq.JSONField
}

type PM_USER struct {
	sq.TableInfo  `ddl:"primarykey={. cols=user_id}"`
	USER_ID       sq.UUIDField
	USERNAME      sq.StringField `ddl:"unique"`
	EMAIL         sq.StringField `ddl:"unique"`
	NAME          sq.StringField
	PASSWORD_HASH sq.StringField
}

type PM_USER_AUTHZ struct {
	sq.TableInfo          `ddl:"primarykey={. cols=site_id,user_id}"`
	SITE_ID               sq.UUIDField
	USER_ID               sq.UUIDField
	ROLES                 sq.CustomField `ddl:"type=JSON postgres:type=TEXT[]"`
	AUTHZ_ATTRIBUTES      sq.JSONField
	ROLE_AUTHZ_ATTRIBUTES sq.JSONField
}

type PM_USER_AUTHZ_ROLES_TBLIDX struct {
	sq.TableInfo `ddl:"ignore=postgres,mysql primarykey={pm_user_authz_roles_tblidx_pkey cols=site_id,user_id,role}"`
	SITE_ID      sq.UUIDField
	USER_ID      sq.UUIDField
	ROLE         sq.StringField
}

type PM_ROLE struct {
	sq.TableInfo     `ddl:"ignore=postgres,mysql primarykey={. cols=site_id,role}"`
	SITE_ID          sq.UUIDField
	ROLE             sq.StringField
	AUTHZ_ATTRIBUTES sq.JSONField
}

type PM_SESSION struct {
	sq.TableInfo `ddl:"primarykey={. cols=session_hash}"`
	SESSION_HASH sq.BlobField
	SITE_ID      sq.UUIDField
	USER_ID      sq.UUIDField
}
