package main

import (
	"database/sql"
	"flag"
	"log"

	"github.com/bokwoon95/pagemanager"
	"github.com/bokwoon95/pagemanager/internal/tables"
	"github.com/bokwoon95/pagemanager/sq"
	"github.com/bokwoon95/pagemanager/sq/ddl"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	flag.Parse()
	cfg, err := pagemanager.DefaultConfig()
	if err != nil {
		log.Fatal(err)
	}
	pm, err := pagemanager.New(cfg)
	if err != nil {
		log.Fatal(err)
	}
	_ = pm
	// fmt.Printf("%#v\n", cfg)

	// err = ensuretables(cfg)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}

func ensuretables(cfg *pagemanager.Config) error {
	var tbls = []sq.SchemaTable{
		&tables.PM_SITE{},
		&tables.PM_PLUGIN{},
		&tables.PM_HANDLER{},
		&tables.PM_CAPABILITY{},
		&tables.PM_ALLOWED_PLUGIN{},
		&tables.PM_DENIED_PLUGIN{},
		&tables.PM_ALLOWED_HANDLER{},
		&tables.PM_DENIED_HANDLER{},
		&tables.PM_ROLE{},
		&tables.PM_TAG{},
		&tables.PM_ROLE_CAPABILITY{},
		&tables.PM_TAG_CAPABILITY{},
		&tables.PM_TAG_OWNER{},
		&tables.PM_USER{},
		&tables.PM_USER_ROLE{},
		&tables.PM_SESSION{},
		&tables.PM_URL{},
		&tables.PM_URL_ROLE_CAPABILITY{},
		&tables.PM_URL_TAG{},
		&tables.PM_TEMPLATE_DATA{},
	}
	for _, tbl := range tbls {
		err := sq.ReflectTable(tbl, "")
		if err != nil {
			return err
		}
	}
	var driverName string
	switch cfg.DatabaseDialect {
	case "sqlite":
		driverName = "sqlite3"
	case "postgres":
		driverName = "postgres"
	default:
		driverName = "mysql"
	}
	db, err := sql.Open(driverName, cfg.DatabaseURL)
	if err != nil {
		return err
	}
	gotMetadata, err := ddl.NewDatabaseMetadata(cfg.DatabaseDialect, ddl.WithDB(db, nil))
	if err != nil {
		return err
	}
	_ = gotMetadata
	wantMetadata, err := ddl.NewDatabaseMetadata(cfg.DatabaseDialect, ddl.WithTables(tbls...))
	if err != nil {
		return err
	}
	_ = wantMetadata
	// m, err := ddl.Migrate(ddl.DropExtraneous, gotMetadata, ddl.DatabaseMetadata{})
	m, err := ddl.Migrate(ddl.CreateMissing, gotMetadata, wantMetadata)
	if err != nil {
		return err
	}
	_ = m
	// err = m.WriteSQL(os.Stdout)
	// if err != nil {
	// 	return err
	// }
	err = m.Exec(db)
	if err != nil {
		return err
	}
	return nil
}
