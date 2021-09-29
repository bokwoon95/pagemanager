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
	sq.TableInfo `ddl:"primarykey={. cols=site_id,path}"`
	SITE_ID      sq.UUIDField
	PATH         sq.StringField
	PLUGIN       sq.StringField
	HANDLER      sq.StringField
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

type PM_ROLE struct {
	sq.TableInfo `ddl:"primarykey={. cols=site_id,role}"`
	SITE_ID      sq.UUIDField
	ROLE         sq.StringField
}

type PM_ROLE_USER struct {
	sq.TableInfo `ddl:"primarykey={. cols=site_id,role,user_id}"`
	SITE_ID      sq.UUIDField
	ROLE         sq.StringField
	USER_ID      sq.UUIDField
}

type PM_PERMISSION struct {
	sq.TableInfo `ddl:"ignore=postgres,mysql primarykey={. cols=site_id,role,label,action}"`
	SITE_ID      sq.UUIDField
	ROLE         sq.StringField
	LABEL        sq.StringField
	ACTION       sq.StringField
}

type PM_SESSION struct {
	sq.TableInfo `ddl:"primarykey={. cols=session_hash}"`
	SESSION_HASH sq.BlobField
	SITE_ID      sq.UUIDField
	USER_ID      sq.UUIDField
}
