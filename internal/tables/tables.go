package tables

import (
	"github.com/bokwoon95/pagemanager/sq"
)

type PM_SITE struct {
	sq.TableInfo `ddl:"primarykey=site_id unique=domain,subdomain,tilde_prefix"`
	SITE_ID      sq.UUIDField    `ddl:"mysql:type=BINARY(16) sqlite:type=BLOB"`
	DOMAIN       sq.StringField  `ddl:"notnull"`
	SUBDOMAIN    sq.StringField  `ddl:"notnull"`
	TILDE_PREFIX sq.StringField  `ddl:"notnull"`
	IS_PRIMARY   sq.BooleanField `ddl:"unique default=NULL"`
}

type PM_PLUGIN struct {
	sq.TableInfo `ddl:"primarykey=plugin unique=url,version"`
	PLUGIN       sq.StringField
	URL          sq.StringField `ddl:"notnull"`
	VERSION      sq.StringField `ddl:"notnull"`
}

type PM_HANDLER struct {
	sq.TableInfo `ddl:"primarykey=plugin,handler"`
	PLUGIN       sq.StringField `ddl:"references={pm_plugin onupdate=cascade index}"`
	HANDLER      sq.StringField
}

type PM_CAPABILITY struct {
	sq.TableInfo `ddl:"primarykey=plugin,capability"`
	PLUGIN       sq.StringField `ddl:"references={pm_plugin onupdate=cascade index}"`
	CAPABILITY   sq.StringField
}

type PM_ALLOWED_PLUGIN struct {
	sq.TableInfo `ddl:"primarykey=plugin,site_id"`
	PLUGIN       sq.StringField `ddl:"references={pm_plugin onupdate=cascade index}"`
	SITE_ID      sq.UUIDField   `ddl:"references={pm_site onupdate=cascade index}"`
}

type PM_DENIED_PLUGIN struct {
	sq.TableInfo `ddl:"primarykey=plugin,site_id"`
	PLUGIN       sq.StringField `ddl:"references={pm_plugin onupdate=cascade index} "`
	SITE_ID      sq.UUIDField   `ddl:"references={pm_site onupdate=cascade index}"`
}

type PM_ALLOWED_HANDLER struct {
	sq.TableInfo `ddl:"primarykey=plugin,handler,site_id"`
	PLUGIN       sq.StringField
	HANDLER      sq.StringField
	SITE_ID      sq.UUIDField `ddl:"references={pm_site onupdate=cascade index}"`

	_ struct{} `ddl:"foreignkey={plugin,handler references=pm_handler onupdate=cascade index}"`
}

type PM_DENIED_HANDLER struct {
	sq.TableInfo `ddl:"primarykey=plugin,handler,site_id"`
	PLUGIN       sq.StringField
	HANDLER      sq.StringField
	SITE_ID      sq.UUIDField `ddl:"references={pm_site onupdate=cascade index}"`

	_ struct{} `ddl:"foreignkey={plugin,handler references=pm_handler onupdate=cascade index}"`
}

type PM_ROLE struct {
	sq.TableInfo `ddl:"primarykey=site_id,plugin,role"`
	SITE_ID      sq.UUIDField   `ddl:"references={pm_site onupdate=cascade index}"`
	PLUGIN       sq.StringField `ddl:"references={pm_plugin onupdate=cascade index}"`
	ROLE         sq.StringField `ddl:""`
}

type PM_TAG struct {
	sq.TableInfo `ddl:"primarykey=site_id,plugin,tag"`
	SITE_ID      sq.UUIDField   `ddl:"references={pm_site onupdate=cascade index}"`
	PLUGIN       sq.StringField `ddl:"references={pm_plugin onupdate=cascade index}"`
	TAG          sq.StringField
}

type PM_ROLE_CAPABILITY struct {
	sq.TableInfo `ddl:"primarykey=site_id,plugin,role,capability"`
	SITE_ID      sq.UUIDField
	PLUGIN       sq.StringField
	ROLE         sq.StringField
	CAPABILITY   sq.StringField

	_ struct{} `ddl:"foreignkey={site_id,plugin,role references=pm_role onupdate=cascade index}"`
	_ struct{} `ddl:"foreignkey={plugin,capabillity references=pm_capability onupdate=cascade index}"`
}

type PM_TAG_CAPABILITY struct {
	sq.TableInfo `ddl:"primarykey=site_id,plugin,tag,role"`
	SITE_ID      sq.UUIDField
	PLUGIN       sq.StringField
	TAG          sq.StringField
	ROLE         sq.StringField
	CAPABILITY   sq.StringField `ddl:"notnull"`

	_ struct{} `ddl:"foreignkey={site_id,plugin,tag references=pm_tag onupdate=cascade index}"`
	_ struct{} `ddl:"foreignkey={site_id,plugin,role references=pm_role onupdate=cascade index}"`
	_ struct{} `ddl:"foreignkey={plugin,capability references=pm_capability onupdate=cascade index}"`
}

type PM_TAG_OWNER struct {
	sq.TableInfo `ddl:"primarykey=site_id,plugin,tag,role"`
	SITE_ID      sq.UUIDField
	PLUGIN       sq.StringField
	TAG          sq.StringField
	ROLE         sq.StringField

	_ struct{} `ddl:"foreignkey={site_id,plugin,tag references=pm_tag onupdate=cascade index}"`
	_ struct{} `ddl:"foreignkey={site_id,plugin,role references=pm_role onupdate=cascade index}"`
}

type PM_USER struct {
	sq.TableInfo           `ddl:"primarykey=site_id,user_id unique=site_id,username unique=site_id,email"`
	SITE_ID                sq.UUIDField `ddl:"references={pm_site onupdate=cascade index}"`
	USER_ID                sq.UUIDField
	NAME                   sq.StringField
	USERNAME               sq.StringField `ddl:"notnull"`
	EMAIL                  sq.StringField `ddl:"notnull"`
	PASSWORD_HASH          sq.StringField
	RESET_PASSWORD_TOKEN   sq.StringField
	RESET_PASSWORD_SENT_AT sq.TimeField

	_ struct{} `ddl:"foreignkey={site_id,plugin,tag references=pm_tag onupdate=cascade index}"`
	_ struct{} `ddl:"foreignkey={site_id,plugin,role references=pm_role onupdate=cascade index}"`
}

type PM_SESSION struct {
	sq.TableInfo `ddl:"primarykey=site_id,session_hash"`
	SITE_ID      sq.UUIDField
	SESSION_HASH sq.BinaryField
	USER_ID      sq.UUIDField `ddl:"notnull"`

	_ struct{} `ddl:"foreignkey={site_id,user_id references=pm_user onupdate=cascade index}"`
}

type PM_USER_ROLE struct {
	sq.TableInfo `ddl:"primarykey=site_id,user_id,plugin,role"`
	SITE_ID      sq.UUIDField
	USER_ID      sq.UUIDField
	PLUGIN       sq.StringField
	ROLE         sq.StringField

	_ struct{} `ddl:"foreignkey={site_id,user_id references=pm_user onupdate=cascade index}"`
	_ struct{} `ddl:"foreignkey={site_id,plugin,role references=pm_role onupdate=cascade index}"`
}

type PM_URL struct {
	sq.TableInfo `ddl:"primarykey=site_id,urlpath"`
	SITE_ID      sq.UUIDField `ddl:"references={pm_site onupdate=cascade index}"`
	URLPATH      sq.StringField
	PLUGIN       sq.StringField
	HANDLER      sq.StringField
	CONFIG       sq.StringField

	_ struct{} `ddl:"foreignkey={plugin,handler references=pm_plugin onupdate=cascade index}"`
}

type PM_URL_ROLE_CAPABILITY struct {
	sq.TableInfo `ddl:"primarykey=site_id,urlpath,plugin,role,capability"`
	SITE_ID      sq.UUIDField
	URLPATH      sq.StringField
	PLUGIN       sq.StringField
	ROLE         sq.StringField
	CAPABILITY   sq.StringField

	_ struct{} `ddl:"foreignkey={site_id,plugin,role references=pm_role onupdate=cascade index}"`
	_ struct{} `ddl:"foreignkey={plugin,capability references=pm_capability onupdate=cascade index}"`
}

type PM_URL_TAG struct {
	sq.TableInfo `ddl:"primarykey=site_id,urlpath,plugin,tag"`
	SITE_ID      sq.UUIDField
	URLPATH      sq.StringField
	PLUGIN       sq.StringField
	TAG          sq.StringField

	_ struct{} `ddl:"foreignkey={site_id,plugin,tag references=pm_tag onupdate=cascade index}"`
}

type PM_TEMPLATE_DATA struct {
	sq.TableInfo `ddl:"primarykey=site_id,langcode,datafile"`
	SITE_ID      sq.UUIDField `ddl:"references={pm_site onupdate=cascade index}"`
	LANGCODE     sq.StringField
	DATA_FILE    sq.StringField
	DATA         sq.JSONField
}
