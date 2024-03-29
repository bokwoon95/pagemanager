package ddl

import (
	"testing"

	"github.com/bokwoon95/pagemanager/sq"
	"github.com/bokwoon95/pagemanager/sq/internal/testutil"
)

func Test_Function(t *testing.T) {
	type TT struct {
		dialect            string
		item               Function
		wantFunctionSchema string
		wantFunctionName   string
		wantRawArgs        string
		wantReturnType     string
	}

	assert := func(t *testing.T, tt TT) {
		err := tt.item.populateFunctionInfo(tt.dialect)
		if err != nil {
			t.Fatal(testutil.Callers(), err)
		}
		sql, _, _, err := sq.ToSQL(tt.dialect, CreateFunctionCommand{Function: tt.item})
		if err != nil {
			t.Fatal(testutil.Callers(), err)
		}
		if diff := testutil.Diff(sql, tt.item.SQL); diff != "" {
			t.Error(testutil.Callers(), diff)
		}
		if diff := testutil.Diff(tt.item.FunctionSchema, tt.wantFunctionSchema); diff != "" {
			t.Error(testutil.Callers(), diff)
		}
		if diff := testutil.Diff(tt.item.FunctionName, tt.wantFunctionName); diff != "" {
			t.Error(testutil.Callers(), diff)
		}
		if diff := testutil.Diff(tt.item.RawArgs, tt.wantRawArgs); diff != "" {
			t.Error(testutil.Callers(), diff)
		}
		if diff := testutil.Diff(tt.item.ReturnType, tt.wantReturnType); diff != "" {
			t.Error(testutil.Callers(), diff)
		}
	}

	t.Run("(dialect == postgres)", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = sq.DialectPostgres
		tt.item.SQL = `CREATE FUNCTION one() RETURNS integer AS $$ lorem ipsum $$ LANGUAGE plpgsql`
		tt.wantFunctionName = "one"
		tt.wantReturnType = "integer"
		assert(t, tt)
	})

	t.Run("(dialect == postgres)", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = sq.DialectPostgres
		tt.item.SQL = `CREATE FUNCTION one(     ) RETURNS integer`
		tt.wantFunctionName = "one"
		tt.wantReturnType = "integer"
		assert(t, tt)
	})

	t.Run("(dialect == postgres)", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = sq.DialectPostgres
		tt.item.SQL = `CREATE FUNCTION app.tf1 (integer, in numeric = 3.14) RETURNS integer LANGUAGE sql AS $$ lorem ipsum $$`
		tt.wantFunctionSchema = "app"
		tt.wantFunctionName = "tf1"
		tt.wantRawArgs = "integer, in numeric = 3.14"
		tt.wantReturnType = "integer"
		assert(t, tt)
	})

	t.Run("(dialect == postgres)", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = sq.DialectPostgres
		tt.item.SQL = `CREATE OR REPLACE FUNCTION double_salary(emp) RETURNS numeric AS $$ lorem ipsum $$ LANGUAGE sql`
		tt.wantFunctionName = "double_salary"
		tt.wantRawArgs = "emp"
		tt.wantReturnType = "numeric"
		assert(t, tt)
	})

	t.Run("(dialect == postgres)", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = sq.DialectPostgres
		tt.item.SQL = `CREATE FUNCTION sum_n_product (int, y int =25, OUT sum int DEFAULT 3, OUT product int=22)`
		tt.wantFunctionName = "sum_n_product"
		tt.wantRawArgs = "int, y int =25, OUT sum int DEFAULT 3, OUT product int=22"
		tt.wantReturnType = ""
		assert(t, tt)
	})

	t.Run("(dialect == postgres)", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = sq.DialectPostgres
		tt.item.SQL = `CREATE FUNCTION make_array(anyelement, anyelement) RETURNS anyarray AS LANGUAGE plpgsql $$ lorem ipsum $$`
		tt.wantFunctionName = "make_array"
		tt.wantRawArgs = "anyelement, anyelement"
		tt.wantReturnType = "anyarray"
		assert(t, tt)
	})

	t.Run("(dialect == postgres)", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = sq.DialectPostgres
		tt.item.SQL = `
CREATE OR REPLACE FUNCTION years_compare( IN year1 integer DEFAULT NULL,
                                          year2 IN integer DEFAULT NULL ) RETURNS BOOLEAN AS $$ lorem ipsum $$ language SQL`
		tt.wantFunctionName = "years_compare"
		tt.wantRawArgs = `IN year1 integer DEFAULT NULL,
                                          year2 IN integer DEFAULT NULL`
		tt.wantReturnType = "BOOLEAN"
		assert(t, tt)
	})

	t.Run("(dialect == postgres)", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = sq.DialectPostgres
		tt.item.SQL = `create function foo(bar varchar , baz varchar='qux') returns varchar AS $$ lorem ipsum $$ language sql`
		tt.wantFunctionName = "foo"
		tt.wantRawArgs = "bar varchar , baz varchar='qux'"
		tt.wantReturnType = "varchar"
		assert(t, tt)
	})

	t.Run("(dialect == postgres)", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = sq.DialectPostgres
		tt.item.SQL = `CREATE FUNCTION get_count_of_earners(salary_val IN decimal, alphabets []TEXT='{"a", "b", "c"}', names VARIADIC [][]text = ARRAY[ARRAY['a', 'b'], ARRAY['c', 'd']]) RETURNS integer AS $$ lorem ipsum $$ language plpgsql`
		tt.wantFunctionName = "get_count_of_earners"
		tt.wantRawArgs = `salary_val IN decimal, alphabets []TEXT='{"a", "b", "c"}', names VARIADIC [][]text = ARRAY[ARRAY['a', 'b'], ARRAY['c', 'd']]`
		tt.wantReturnType = "integer"
		assert(t, tt)
	})

	t.Run("(dialect != postgres)", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = sq.DialectMySQL
		tt.item.SQL = `CREATE FUNCTION hello (s CHAR(20)) RETURNS CHAR(50) DETERMINISTIC RETURN CONCAT('Hello, ',s,'!')`
		assert(t, tt)
	})

	t.Run("(dialect == postgres) no opening bracket", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = sq.DialectPostgres
		tt.item.SQL = `CREATE FUNCTION temp`
		err := tt.item.populateFunctionInfo(tt.dialect)
		if err == nil {
			t.Fatal(testutil.Callers(), "expected error but got nil")
		}
	})

	t.Run("(dialect == postgres) no closing bracket", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = sq.DialectPostgres
		tt.item.SQL = `CREATE FUNCTION temp(`
		err := tt.item.populateFunctionInfo(tt.dialect)
		if err == nil {
			t.Fatal(testutil.Callers(), "expected error but got nil")
		}
	})
}

func Test_DropFunctionCommand(t *testing.T) {
	type TT struct {
		dialect   string
		item      Command
		wantQuery string
		wantArgs  []interface{}
	}

	assert := func(t *testing.T, tt TT) {
		gotQuery, gotArgs, _, err := sq.ToSQL(tt.dialect, tt.item)
		if err != nil {
			t.Fatal(testutil.Callers(), err)
		}
		if diff := testutil.Diff(gotQuery, tt.wantQuery); diff != "" {
			t.Error(testutil.Callers(), diff)
		}
		if diff := testutil.Diff(gotArgs, tt.wantArgs); diff != "" {
			t.Error(testutil.Callers(), diff)
		}
	}

	t.Run("(dialect == postgres)", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = sq.DialectPostgres
		tt.item = DropFunctionCommand{
			DropIfExists: true,
			Function: Function{
				FunctionSchema: "public",
				FunctionName:   "my_function",
				RawArgs:        "TEXT, INT",
			},
			DropCascade: true,
		}
		tt.wantQuery = `DROP FUNCTION IF EXISTS public.my_function(TEXT, INT) CASCADE`
		assert(t, tt)
	})

	t.Run("(dialect == sqlite)", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = sq.DialectSQLite
		tt.item = DropFunctionCommand{
			Function: Function{FunctionName: "my_function"},
		}
		_, _, _, err := sq.ToSQL(tt.dialect, tt.item)
		if err == nil {
			t.Fatal(testutil.Callers(), "expected error but got nil")
		}
	})
}
