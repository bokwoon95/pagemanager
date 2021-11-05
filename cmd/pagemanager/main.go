package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/bokwoon95/pagemanager"
	"github.com/bokwoon95/pagemanager/internal/tables"
	"github.com/bokwoon95/pagemanager/sq"
	"github.com/bokwoon95/pagemanager/sq/ddl"
)

func main() {
	flag.Parse()
	cfg, err := pagemanager.DefaultConfig()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%#v\n", cfg)

	var tbls = []sq.SchemaTable{
		&tables.PM_SITE{},
		&tables.PM_PLUGIN{},
		&tables.PM_HANDLER{},
	}
	for _, tbl := range tbls {
		sq.ReflectTable(tbl, "")
	}

	dbm, err := ddl.NewDatabaseMetadata(sq.DialectSQLite, ddl.WithTables(tbls...))
	if err != nil {
		log.Fatal(err)
	}
	m, err := ddl.Migrate(ddl.CreateMissing|ddl.UpdateExisting, ddl.DatabaseMetadata{}, dbm)
	if err != nil {
		log.Fatal(err)
	}
	m.WriteSQL(os.Stdout)
}
