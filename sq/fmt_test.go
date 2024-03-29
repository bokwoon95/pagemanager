package sq

import (
	"bytes"
	"database/sql"
	"testing"

	"github.com/bokwoon95/pagemanager/sq/internal/testutil"
)

func Test_Fprintf(t *testing.T) {
	USERS := struct {
		tmptable
		USER_ID tmpfield
		NAME    tmpfield
		EMAIL   tmpfield
		AGE     tmpfield
	}{
		tmptable: [2]string{"", "users"},
		USER_ID:  [2]string{"", "user_id"},
		NAME:     [2]string{"", "name"},
		EMAIL:    [2]string{"", "email"},
		AGE:      [2]string{"", "age"},
	}

	type TT struct {
		dialect    string
		format     string
		values     []interface{}
		wantQuery  string
		wantArgs   []interface{}
		wantParams map[string][]int
	}

	assert := func(t *testing.T, tt TT) {
		buf := bufpool.Get().(*bytes.Buffer)
		defer func() {
			buf.Reset()
			bufpool.Put(buf)
		}()
		gotArgs, gotParams := []interface{}{}, map[string][]int{}
		err := BufferPrintf(tt.dialect, buf, &gotArgs, gotParams, nil, nil, tt.format, tt.values)
		if err != nil {
			t.Fatal(testutil.Callers(), err)
		}
		if diff := testutil.Diff(buf.String(), tt.wantQuery); diff != "" {
			t.Error(testutil.Callers(), diff)
		}
		if diff := testutil.Diff(gotArgs, tt.wantArgs); diff != "" {
			t.Error(testutil.Callers(), diff)
		}
		if diff := testutil.Diff(gotParams, tt.wantParams); diff != "" {
			t.Error(testutil.Callers(), diff)
		}
	}

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.wantArgs = []interface{}{}
		tt.wantParams = map[string][]int{}
		assert(t, tt)
	})

	t.Run("escape curly bracket", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.format = "SELECT {} = '{{}'"
		tt.values = []interface{}{"{}"}
		tt.wantQuery = "SELECT ? = '{}'"
		tt.wantArgs = []interface{}{"{}"}
		tt.wantParams = map[string][]int{}
		assert(t, tt)
	})

	t.Run("expr", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.format = "(MAX(AVG({avg1}), AVG({avg2}), SUM({sum})) + {incr}) IN ({slice})"
		tt.values = []interface{}{
			Param("avg1", USERS.USER_ID),
			Param("avg2", USERS.AGE),
			Param("sum", USERS.AGE),
			Param("incr", 1),
			Param("slice", []int{1, 2, 3}),
		}
		tt.wantQuery = "(MAX(AVG(user_id), AVG(age), SUM(age)) + ?) IN (?, ?, ?)"
		tt.wantArgs = []interface{}{1, 1, 2, 3}
		tt.wantParams = map[string][]int{"incr": {0}}
		assert(t, tt)
	})

	t.Run("nameonly", func(t *testing.T) {
		t.Parallel()
		USERS := struct {
			tmptable
			USER_ID tmpfield
			NAME    tmpfield
			EMAIL   tmpfield
			AGE     tmpfield
		}{
			tmptable: [2]string{"public", "users"},
			USER_ID:  [2]string{"users", "user_id"},
			NAME:     [2]string{"users", "name"},
			EMAIL:    [2]string{"users", "email"},
			AGE:      [2]string{"users", "age"},
		}
		var tt TT
		tt.format = "UPDATE {table}, {table:nameonly} SET {name} = {3:nameonly} WHERE {user_id} = NEW.{user_id:nameonly}"
		tt.values = []interface{}{
			Param("table", USERS),
			Param("name", USERS.NAME),
			"bob",
			Param("user_id", USERS.USER_ID),
		}
		tt.wantQuery = "UPDATE public.users, users SET users.name = ? WHERE users.user_id = NEW.user_id"
		tt.wantArgs = []interface{}{"bob"}
		tt.wantParams = map[string][]int{}
		assert(t, tt)
	})

	t.Run("mysql anonymous", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectMySQL
		tt.format = "SELECT {} FROM {} WHERE {} = {} AND {} <> {} AND {} IN ({})"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, 5, USERS.EMAIL, "bob@email.com", USERS.NAME, []string{"tom", "dick", "harry"}}
		tt.wantQuery = "SELECT name FROM users WHERE age = ? AND email <> ? AND name IN (?, ?, ?)"
		tt.wantArgs = []interface{}{5, "bob@email.com", "tom", "dick", "harry"}
		tt.wantParams = map[string][]int{}
		assert(t, tt)
	})

	t.Run("mysql ordinal", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectMySQL
		tt.format = "SELECT {} FROM {} WHERE {} = {5} AND {} <> {5} AND {1} IN ({6})"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, USERS.EMAIL, "bob@email.com", []string{"tom", "dick", "harry"}}
		tt.wantQuery = "SELECT name FROM users WHERE age = ? AND email <> ? AND name IN (?, ?, ?)"
		tt.wantArgs = []interface{}{"bob@email.com", "bob@email.com", "tom", "dick", "harry"}
		tt.wantParams = map[string][]int{}
		assert(t, tt)
	})

	t.Run("mysql Param", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectMySQL
		tt.format = "SELECT {} FROM {} WHERE {3} = {age} AND {3} > {age} AND {4} <> {email} AND {1} IN ({names})"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, USERS.EMAIL, Param("email", "bob@email.com"), Param("age", 5), Param("names", []string{"tom", "dick", "harry"})}
		tt.wantQuery = "SELECT name FROM users WHERE age = ? AND age > ? AND email <> ? AND name IN (?, ?, ?)"
		tt.wantArgs = []interface{}{5, 5, "bob@email.com", "tom", "dick", "harry"}
		tt.wantParams = map[string][]int{"age": {0, 1}, "email": {2}}
		assert(t, tt)
	})

	t.Run("mysql sql.Named", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectMySQL
		tt.format = "SELECT {} FROM {} WHERE {3} = {age} AND {3} > {age} AND {4} <> {email}"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, USERS.EMAIL, sql.Named("email", "bob@email.com"), sql.Named("age", 5)}
		err := BufferPrintf(tt.dialect, new(bytes.Buffer), new([]interface{}), make(map[string][]int), nil, nil, tt.format, tt.values)
		if err == nil {
			t.Error(testutil.Callers(), "expected error but got nil")
		}
	})

	t.Run("postgres anonymous", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectPostgres
		tt.format = "SELECT {} FROM {} WHERE {} = {} AND {} <> {} AND {} IN ({})"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, 5, USERS.EMAIL, "bob@email.com", USERS.NAME, []string{"tom", "dick", "harry"}}
		tt.wantQuery = "SELECT name FROM users WHERE age = $1 AND email <> $2 AND name IN ($3, $4, $5)"
		tt.wantArgs = []interface{}{5, "bob@email.com", "tom", "dick", "harry"}
		tt.wantParams = map[string][]int{}
		assert(t, tt)
	})

	t.Run("postgres ordinal", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectPostgres
		tt.format = "SELECT {} FROM {} WHERE {} = {5} AND {} <> {5} AND {1} IN ({6}) AND {4} IN ({6})"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, USERS.EMAIL, "bob@email.com", []string{"tom", "dick", "harry"}}
		tt.wantQuery = "SELECT name FROM users WHERE age = $1 AND email <> $1 AND name IN ($2, $3, $4) AND email IN ($5, $6, $7)"
		tt.wantArgs = []interface{}{"bob@email.com", "tom", "dick", "harry", "tom", "dick", "harry"}
		tt.wantParams = map[string][]int{}
		assert(t, tt)
	})

	t.Run("postgres Param", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectPostgres
		tt.format = "SELECT {} FROM {} WHERE {3} = {age} AND {3} > {age} AND {4} <> {email} AND {1} IN ({names}) AND {4} IN ({names})"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, USERS.EMAIL, Param("email", "bob@email.com"), Param("age", 5), Param("names", []string{"tom", "dick", "harry"})}
		tt.wantQuery = "SELECT name FROM users WHERE age = $1 AND age > $1 AND email <> $2 AND name IN ($3, $4, $5) AND email IN ($6, $7, $8)"
		tt.wantArgs = []interface{}{5, "bob@email.com", "tom", "dick", "harry", "tom", "dick", "harry"}
		tt.wantParams = map[string][]int{"age": {0}, "email": {1}}
		assert(t, tt)
	})

	t.Run("postgres sql.Named", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectPostgres
		tt.format = "SELECT {} FROM {} WHERE {3} = {age} AND {3} > {age} AND {4} <> {email}"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, USERS.EMAIL, sql.Named("email", "bob@email.com"), sql.Named("age", 5)}
		err := BufferPrintf(tt.dialect, new(bytes.Buffer), new([]interface{}), make(map[string][]int), nil, nil, tt.format, tt.values)
		if err == nil {
			t.Error(testutil.Callers(), "expected error but got nil")
		}
	})

	t.Run("sqlite anonymous", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLite
		tt.format = "SELECT {} FROM {} WHERE {} = {} AND {} <> {} AND {} IN ({})"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, 5, USERS.EMAIL, "bob@email.com", USERS.NAME, []string{"tom", "dick", "harry"}}
		tt.wantQuery = "SELECT name FROM users WHERE age = $1 AND email <> $2 AND name IN ($3, $4, $5)"
		tt.wantArgs = []interface{}{5, "bob@email.com", "tom", "dick", "harry"}
		tt.wantParams = map[string][]int{}
		assert(t, tt)
	})

	t.Run("sqlite ordinal", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLite
		tt.format = "SELECT {} FROM {} WHERE {} = {5} AND {} <> {5} AND {1} IN ({6}) AND {4} IN ({6})"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, USERS.EMAIL, "bob@email.com", []string{"tom", "dick", "harry"}}
		tt.wantQuery = "SELECT name FROM users WHERE age = $1 AND email <> $1 AND name IN ($2, $3, $4) AND email IN ($5, $6, $7)"
		tt.wantArgs = []interface{}{"bob@email.com", "tom", "dick", "harry", "tom", "dick", "harry"}
		tt.wantParams = map[string][]int{}
		assert(t, tt)
	})

	t.Run("sqlite Param", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLite
		tt.format = "SELECT {} FROM {} WHERE {3} = {age} AND {3} > {age} AND {4} <> {email} AND {1} IN ({names}) AND {4} IN ({names})"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, USERS.EMAIL, Param("email", "bob@email.com"), Param("age", 5), Param("names", []string{"tom", "dick", "harry"})}
		tt.wantQuery = "SELECT name FROM users WHERE age = $1 AND age > $1 AND email <> $2 AND name IN ($3, $4, $5) AND email IN ($6, $7, $8)"
		tt.wantArgs = []interface{}{5, "bob@email.com", "tom", "dick", "harry", "tom", "dick", "harry"}
		tt.wantParams = map[string][]int{"age": {0}, "email": {1}}
		assert(t, tt)
	})

	t.Run("sqlite sql.Named", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLite
		tt.format = "SELECT {} FROM {} WHERE {3} = {age} AND {3} > {age} AND {4} <> {email}"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, USERS.EMAIL, sql.Named("email", "bob@email.com"), sql.Named("age", 5)}
		tt.wantQuery = "SELECT name FROM users WHERE age = $age AND age > $age AND email <> $email"
		tt.wantArgs = []interface{}{sql.Named("age", 5), sql.Named("email", "bob@email.com")}
		tt.wantParams = map[string][]int{"age": {0}, "email": {1}}
		assert(t, tt)
	})

	t.Run("sqlite mixing sql.Named and sq.Param", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLite
		tt.format = "SELECT {} FROM {} WHERE {3} = {age} AND {3} > {age} AND {4} <> {email}"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, USERS.EMAIL, Param("email", "bob@email.com"), sql.Named("age", 5)}
		tt.wantQuery = "SELECT name FROM users WHERE age = $age AND age > $age AND email <> $2"
		tt.wantArgs = []interface{}{sql.Named("age", 5), "bob@email.com"}
		tt.wantParams = map[string][]int{"age": {0}, "email": {1}}
		assert(t, tt)
	})

	t.Run("sqlite mixing sql.Named and sq.Param", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLite
		tt.format = "SELECT {} FROM {} WHERE {4} <> {email} AND {3} = {age} AND {3} > {age}"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, USERS.EMAIL, sql.Named("age", 5), Param("email", "bob@email.com")}
		tt.wantQuery = "SELECT name FROM users WHERE email <> $1 AND age = $age AND age > $age"
		tt.wantArgs = []interface{}{"bob@email.com", sql.Named("age", 5)}
		tt.wantParams = map[string][]int{"age": {1}, "email": {0}}
		assert(t, tt)
	})

	t.Run("mssql anonymous", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLServer
		tt.format = "SELECT {} FROM {} WHERE {} = {} AND {} <> {} AND {} IN ({})"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, 5, USERS.EMAIL, "bob@email.com", USERS.NAME, []string{"tom", "dick", "harry"}}
		tt.wantQuery = "SELECT name FROM users WHERE age = @p1 AND email <> @p2 AND name IN (@p3, @p4, @p5)"
		tt.wantArgs = []interface{}{5, "bob@email.com", "tom", "dick", "harry"}
		tt.wantParams = map[string][]int{}
		assert(t, tt)
	})

	t.Run("mssql ordinal", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLServer
		tt.format = "SELECT {} FROM {} WHERE {} = {5} AND {} <> {5} AND {1} IN ({6}) AND {4} IN ({6})"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, USERS.EMAIL, "bob@email.com", []string{"tom", "dick", "harry"}}
		tt.wantQuery = "SELECT name FROM users WHERE age = @p1 AND email <> @p1 AND name IN (@p2, @p3, @p4) AND email IN (@p5, @p6, @p7)"
		tt.wantArgs = []interface{}{"bob@email.com", "tom", "dick", "harry", "tom", "dick", "harry"}
		tt.wantParams = map[string][]int{}
		assert(t, tt)
	})

	t.Run("MSSQL Param", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLServer
		tt.format = "SELECT {} FROM {} WHERE {3} = {age} AND {3} > {age} AND {4} <> {email} AND {1} IN ({names}) AND {4} IN ({names})"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, USERS.EMAIL, Param("email", "bob@email.com"), Param("age", 5), Param("names", []string{"tom", "dick", "harry"})}
		tt.wantQuery = "SELECT name FROM users WHERE age = @p1 AND age > @p1 AND email <> @p2 AND name IN (@p3, @p4, @p5) AND email IN (@p6, @p7, @p8)"
		tt.wantArgs = []interface{}{5, "bob@email.com", "tom", "dick", "harry", "tom", "dick", "harry"}
		tt.wantParams = map[string][]int{"age": {0}, "email": {1}}
		assert(t, tt)
	})

	t.Run("MSSQL sql.Named", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLServer
		tt.format = "SELECT {} FROM {} WHERE {3} = {age} AND {3} > {age} AND {4} <> {email}"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, USERS.EMAIL, sql.Named("email", "bob@email.com"), sql.Named("age", 5)}
		tt.wantQuery = "SELECT name FROM users WHERE age = @age AND age > @age AND email <> @email"
		tt.wantArgs = []interface{}{sql.Named("age", 5), sql.Named("email", "bob@email.com")}
		tt.wantParams = map[string][]int{"age": {0}, "email": {1}}
		assert(t, tt)
	})

	t.Run("MSSQL mixing sql.Named and sq.Param", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLServer
		tt.format = "SELECT {} FROM {} WHERE {3} = {age} AND {3} > {age} AND {4} <> {email}"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, USERS.EMAIL, Param("email", "bob@email.com"), sql.Named("age", 5)}
		tt.wantQuery = "SELECT name FROM users WHERE age = @age AND age > @age AND email <> @p2"
		tt.wantArgs = []interface{}{sql.Named("age", 5), "bob@email.com"}
		tt.wantParams = map[string][]int{"age": {0}, "email": {1}}
		assert(t, tt)
	})

	t.Run("MSSQL mixing sql.Named and sq.Param", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLServer
		tt.format = "SELECT {} FROM {} WHERE {4} <> {email} AND {3} = {age} AND {3} > {age}"
		tt.values = []interface{}{USERS.NAME, USERS, USERS.AGE, USERS.EMAIL, sql.Named("age", 5), Param("email", "bob@email.com")}
		tt.wantQuery = "SELECT name FROM users WHERE email <> @p1 AND age = @age AND age > @age"
		tt.wantArgs = []interface{}{"bob@email.com", sql.Named("age", 5)}
		tt.wantParams = map[string][]int{"age": {1}, "email": {0}}
		assert(t, tt)
	})
}

func Test_Sprintf(t *testing.T) {
	type TT struct {
		dialect    string
		query      string
		args       []interface{}
		wantString string
	}

	assert := func(t *testing.T, tt TT) {
		gotString, err := Sprintf(tt.dialect, tt.query, tt.args)
		if err != nil {
			t.Fatal(testutil.Callers(), err)
		}
		if diff := testutil.Diff(gotString, tt.wantString); diff != "" {
			t.Error(testutil.Callers(), diff)
		}
	}

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = ""
		tt.query = ""
		tt.args = []interface{}{}
		tt.wantString = ""
		assert(t, tt)
	})

	t.Run("insideString, insideIdentifier and escaping single quotes", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = ""
		tt.query = `SELECT ?` +
			`, 'do not "rebind" ? ? ?'` + // string
			`, "do not 'rebind' ? ? ?"` + // identifier
			`, ?` +
			`, ?`
		tt.args = []interface{}{
			"normal string",
			"string with 'quotes' must be escaped",
			"string with already escaped ''quotes'' except for 'this'",
		}
		tt.wantString = `SELECT 'normal string'` +
			`, 'do not "rebind" ? ? ?'` +
			`, "do not 'rebind' ? ? ?"` +
			`, 'string with ''quotes'' must be escaped'` +
			`, 'string with already escaped ''quotes'' except for ''this'''`
		assert(t, tt)
	})

	t.Run("insideString, insideIdentifier and escaping single quotes (dialect == mysql)", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectMySQL
		tt.query = `SELECT ?` +
			`, 'do not "rebind" ? ? ?'` + // string
			", `do not \" 'rebind' ? ? ?`" + // identifier
			", \"do not ``` 'rebind' ? ? ?\"" + // identifier
			`, ?` +
			`, ?`
		tt.args = []interface{}{
			"normal string",
			"string with 'quotes' must be escaped",
			"string with already escaped ''quotes'' except for 'this'",
		}
		tt.wantString = `SELECT 'normal string'` +
			`, 'do not "rebind" ? ? ?'` +
			", `do not \" 'rebind' ? ? ?`" +
			", \"do not ``` 'rebind' ? ? ?\"" +
			`, 'string with ''quotes'' must be escaped'` +
			`, 'string with already escaped ''quotes'' except for ''this'''`
		assert(t, tt)
	})

	t.Run("insideString, insideIdentifier and escaping single quotes (dialect == sqlserver)", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLServer
		tt.query = `SELECT ?` +
			`, 'do not [[rebind] @p1 @p2 @name'` + // string
			", [do not \" 'rebind' [[[[[@pp]] @p3 @p1]" + // identifier
			", \"do not [[[ 'rebind' [[[[[@pp]] @p3 @p1\"" + // identifier
			`, ?` +
			`, @p3`
		tt.args = []interface{}{
			"normal string",
			"string with 'quotes' must be escaped",
			"string with already escaped ''quotes'' except for 'this'",
		}
		tt.wantString = `SELECT 'normal string'` +
			`, 'do not [[rebind] @p1 @p2 @name'` +
			", [do not \" 'rebind' [[[[[@pp]] @p3 @p1]" +
			", \"do not [[[ 'rebind' [[[[[@pp]] @p3 @p1\"" +
			`, 'string with ''quotes'' must be escaped'` +
			`, 'string with already escaped ''quotes'' except for ''this'''`
		assert(t, tt)
	})

	t.Run("mysql", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectMySQL
		tt.query = "SELECT name FROM users WHERE age = ? AND email <> ? AND name IN (?, ?, ?)"
		tt.args = []interface{}{5, "bob@email.com", "tom", "dick", "harry"}
		tt.wantString = "SELECT name FROM users WHERE age = 5 AND email <> 'bob@email.com' AND name IN ('tom', 'dick', 'harry')"
		assert(t, tt)
	})

	t.Run("mysql insideString", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectMySQL
		tt.query = "SELECT name FROM users WHERE age = ? AND email <> '? ? ? ? ''bruh ?' AND name IN (?, ?) ?"
		tt.args = []interface{}{5, "tom", "dick", "harry"}
		tt.wantString = "SELECT name FROM users WHERE age = 5 AND email <> '? ? ? ? ''bruh ?' AND name IN ('tom', 'dick') 'harry'"
		assert(t, tt)
	})

	t.Run("omitted dialect insideString", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = ""
		tt.query = "SELECT name FROM users WHERE age = ? AND email <> '? ? ? ? ''bruh ?' AND name IN (?, ?) ?"
		tt.args = []interface{}{5, "tom", "dick", "harry"}
		tt.wantString = "SELECT name FROM users WHERE age = 5 AND email <> '? ? ? ? ''bruh ?' AND name IN ('tom', 'dick') 'harry'"
		assert(t, tt)
	})

	t.Run("postgres", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectPostgres
		tt.query = "SELECT name FROM users WHERE age = $1 AND email <> $2 AND name IN ($2, $3, $4, $1)"
		tt.args = []interface{}{5, "tom", "dick", "harry"}
		tt.wantString = "SELECT name FROM users WHERE age = 5 AND email <> 'tom' AND name IN ('tom', 'dick', 'harry', 5)"
		assert(t, tt)
	})

	t.Run("postgres insideString", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectPostgres
		tt.query = "SELECT name FROM users WHERE age = $1 AND email <> '$2 $2 $3 $4 ''bruh $1' AND name IN ($2, $3) $4"
		tt.args = []interface{}{5, "tom", "dick", "harry"}
		tt.wantString = "SELECT name FROM users WHERE age = 5 AND email <> '$2 $2 $3 $4 ''bruh $1' AND name IN ('tom', 'dick') 'harry'"
		assert(t, tt)
	})

	t.Run("sqlite", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLite
		tt.query = "SELECT name FROM users WHERE age = $1 AND email <> $2 AND name IN ($2, $3, $4, $1)"
		tt.args = []interface{}{5, "tom", "dick", "harry"}
		tt.wantString = "SELECT name FROM users WHERE age = 5 AND email <> 'tom' AND name IN ('tom', 'dick', 'harry', 5)"
		assert(t, tt)
	})

	t.Run("sqlite insideString", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLite
		tt.query = "SELECT name FROM users WHERE age = $1 AND email <> '$2 $2 $3 $4 ''bruh $1' AND name IN ($2, $3) $4"
		tt.args = []interface{}{5, "tom", "dick", "harry"}
		tt.wantString = "SELECT name FROM users WHERE age = 5 AND email <> '$2 $2 $3 $4 ''bruh $1' AND name IN ('tom', 'dick') 'harry'"
		assert(t, tt)
	})

	t.Run("sqlite mixing ordinal param and named param", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLite
		tt.query = "SELECT name FROM users WHERE age = $age AND age > $1 AND email <> $email"
		tt.args = []interface{}{sql.Named("age", 5), sql.Named("email", "bob@email.com")}
		tt.wantString = "SELECT name FROM users WHERE age = 5 AND age > 5 AND email <> 'bob@email.com'"
		assert(t, tt)
	})

	t.Run("sqlite supports everything", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLite
		tt.query = "SELECT name FROM users WHERE age = ?age AND email <> :email AND name IN (@3, ?4, $5, :5) ? ?"
		tt.args = []interface{}{sql.Named("age", 5), sql.Named("email", "bob@email.com"), "tom", "dick", "harry"}
		tt.wantString = "SELECT name FROM users WHERE age = 5 AND email <> 'bob@email.com' AND name IN ('tom', 'dick', 'harry', 'harry') 5 'bob@email.com'"
		assert(t, tt)
	})

	t.Run("sqlserver", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLServer
		tt.query = "SELECT name FROM users WHERE age = @p1 AND email <> @P2 AND name IN (@p2, @p3, @p4, @P1)"
		tt.args = []interface{}{5, "tom", "dick", "harry"}
		tt.wantString = "SELECT name FROM users WHERE age = 5 AND email <> 'tom' AND name IN ('tom', 'dick', 'harry', 5)"
		assert(t, tt)
	})

	t.Run("sqlserver insideString", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLServer
		tt.query = "SELECT name FROM users WHERE age = @p1 AND email <> '@p2 @p2 @p3 @p4 ''bruh @p1' AND name IN (@p2, @p3) @p4"
		tt.args = []interface{}{5, "tom", "dick", "harry"}
		tt.wantString = "SELECT name FROM users WHERE age = 5 AND email <> '@p2 @p2 @p3 @p4 ''bruh @p1' AND name IN ('tom', 'dick') 'harry'"
		assert(t, tt)
	})

	t.Run("sqlserver mixing ordinal param and named param", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.dialect = DialectSQLServer
		tt.query = "SELECT name FROM users WHERE age = @age AND age > @p1 AND email <> @email"
		tt.args = []interface{}{sql.Named("age", 5), sql.Named("email", "bob@email.com")}
		tt.wantString = "SELECT name FROM users WHERE age = 5 AND age > 5 AND email <> 'bob@email.com'"
		assert(t, tt)
	})
}
