package pagemanager

import (
	sq "github.com/bokwoon95/sq2"
)

type PM_SITE struct {
	sq.TableInfo `ddl:"unique={. cols=domain,subdomain,tilde_prefix}"`
	SITE_ID      sq.UUIDField    `ddl:"notnull primarykey"`
	DOMAIN       sq.StringField  `ddl:"notnull postgres:collate=C"`
	SUBDOMAIN    sq.StringField  `ddl:"notnull postgres:collate=C"`
	TILDE_PREFIX sq.StringField  `ddl:"notnull postgres:collate=C"`
	IS_PRIMARY   sq.BooleanField `ddl:"unique default=NULL"`
}

type PM_PLUGIN struct {
	sq.TableInfo `ddl:"unique={. cols=url,version}"`
	PLUGIN       sq.StringField `ddl:"notnull primarykey postgres:collate=C"`
	URL          sq.StringField `ddl:"notnull postgres:collate=C"`
	VERSION      sq.StringField `ddl:"notnull postgres:collate=C"`
}

type PM_HANDLER struct {
	sq.TableInfo `ddl:"primarykey={. cols=plugin,handler}"`
	PLUGIN       sq.StringField `ddl:"notnull postgres:collate=C references={pm_plugin.plugin onupdate=cascade} index"`
	HANDLER      sq.StringField `ddl:"notnull postgres:collate=C"`
}

type PM_CAPABILITY struct {
	sq.TableInfo `ddl:"primarykey={. cols=plugin,capability}"`
	PLUGIN       sq.StringField `ddl:"notnull postgres:collate=C references={pm_plugin.plugin onupdate=cascade} index"`
	CAPABILITY   sq.StringField `ddl:"notnull postgres:collate=C"`
}

type PM_ALLOWED_PLUGIN struct {
	sq.TableInfo `ddl:"primarykey={. cols=plugin,site_id}"`
	PLUGIN       sq.StringField `ddl:"notnull postgres:collate=C references={pm_plugin.plugin onupdate=cascade} index"`
	SITE_ID      sq.UUIDField   `ddl:"notnull references={pm_site.site_id onupdate=cascade} index"`
}

type PM_DENIED_PLUGIN struct {
	sq.TableInfo `ddl:"primarykey={. cols=plugin,site_id}"`
	PLUGIN       sq.StringField `ddl:"notnull postgres:collate=C references={pm_plugin.plugin onupdate=cascade} index"`
	SITE_ID      sq.UUIDField   `ddl:"notnull references={pm_site.site_id onupdate=cascade} index"`
}

type PM_ALLOWED_HANDLER struct {
	sq.TableInfo `ddl:"primarykey={. cols=plugin,handler,site_id} references={pm_handler.plugin,handler cols=plugin,handler onupdate=cascade} index"`
	PLUGIN       sq.StringField `ddl:"notnull postgres:collate=C"`
	HANDLER      sq.StringField `ddl:"notnull postgres:collate=C"`
	SITE_ID      sq.UUIDField   `ddl:"notnull references={pm_site.site_id onupdate=cascade} index"`
}

type PM_DENIED_HANDLER struct {
	sq.TableInfo `ddl:"primarykey={. cols=plugin,handler,site_id} references={pm_handler.plugin,handler cols=plugin,handler onupdate=cascade} index"`

	PLUGIN  sq.StringField `ddl:"notnull postgres:collate=C"`
	HANDLER sq.StringField `ddl:"notnull postgres:collate=C"`
	SITE_ID sq.UUIDField   `ddl:"notnull references={pm_site.site_id onupdate=cascade} index"`
	_       struct{}       `ddl:""` // TODO: facilitate this hacky workaround (in ddl) to stuff more annotations into a struct
	_       struct{}       // anything that is not a field would get interpeted as a table-level constraint instead
	_       struct{}       // so you can use virtual=fts5 in a non-field as well
	_       struct{}       `ddl:"foreignkey={plugin,handler references=pm_handler.plugin,handler onupdate=cascade}"` // TODO
}

type PM_ROLE struct {
	sq.TableInfo `ddl:"primarykey={. cols=plugin,handler,site_id} references={pm_handler.plugin,handler cols=plugin,handler onupdate=cascade} index"`
}