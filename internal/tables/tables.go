package tables

import (
	sq "github.com/bokwoon95/sq2"
)

type PM_SITE struct {
	sq.TableInfo `ddl:"unique={. cols=domain,subdomain,tilde_prefix}"`
	SITE_ID      sq.UUIDField    `ddl:"notnull primarykey mysql:type=BINARY(16) sqlite:type=BLOB"`
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
	sq.TableInfo `ddl:"primarykey=plugin,site_id"`
	PLUGIN       sq.StringField `ddl:"notnull postgres:collate=C references={pm_plugin onupdate=cascade index}"`
	SITE_ID      sq.UUIDField   `ddl:"notnull references={pm_site onupdate=cascade index}"`
}

type PM_DENIED_PLUGIN struct {
	sq.TableInfo `ddl:"primarykey={. cols=plugin,site_id}"`
	PLUGIN       sq.StringField `ddl:"notnull postgres:collate=C references={pm_plugin.plugin onupdate=cascade} index"`
	SITE_ID      sq.UUIDField   `ddl:"notnull references={pm_site.site_id onupdate=cascade} index"`
}

type PM_ALLOWED_HANDLER struct {
	sq.TableInfo `ddl:"primarykey=plugin,handler,site_id"`
	PLUGIN       sq.StringField `ddl:"notnull postgres:collate=C"`
	HANDLER      sq.StringField `ddl:"notnull postgres:collate=C"`
	SITE_ID      sq.UUIDField   `ddl:"notnull references={pm_site onupdate=cascade} index"`

	// UH OH TODO: primarykey,index and unique's first args should be the comma separated cols, not the name
	// name should be a submodifier
	// wow how interesting, if you put an 'index' submodifier inside a foreignkey or references constraint it will automatically add an index for that foreign key.
	_ struct{} `ddl:"foreignkey={plugin,handler references=pm_handler onupdate=cascade index}"`
	_ struct{} `ddl:"index=plugin,handler"`
}

type PM_DENIED_HANDLER struct {
	sq.TableInfo `ddl:"primarykey={. cols=plugin,handler,site_id} references={pm_handler.plugin,handler cols=plugin,handler onupdate=cascade} index"`

	PLUGIN  sq.StringField `ddl:"notnull postgres:collate=C"`
	HANDLER sq.StringField `ddl:"notnull postgres:collate=C"`
	SITE_ID sq.UUIDField   `ddl:"notnull references={pm_site.site_id onupdate=cascade} index"`

	_ struct{} `ddl:""`
	_ struct{}
	_ struct{} // so you can use virtual=fts5 in a non-field as well
	_ struct{} `ddl:"foreignkey={plugin,handler references=pm_handler.plugin,handler onupdate=cascade}"` // TODO
}

type PM_ROLE struct {
	sq.TableInfo `ddl:"primarykey={. cols=plugin,handler,site_id} references={pm_handler.plugin,handler cols=plugin,handler onupdate=cascade} index"`
}
