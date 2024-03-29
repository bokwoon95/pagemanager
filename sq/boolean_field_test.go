package sq

import (
	"bytes"
	"testing"

	"github.com/bokwoon95/pagemanager/sq/internal/testutil"
)

func Test_BooleanField(t *testing.T) {
	type TT struct {
		dialect                 string
		item                    SQLExcludeAppender
		excludedTableQualifiers []string
		wantQuery               string
		wantArgs                []interface{}
	}

	assert := func(t *testing.T, tt TT) {
		buf := bufpool.Get().(*bytes.Buffer)
		defer func() {
			buf.Reset()
			bufpool.Put(buf)
		}()
		gotArgs, gotParams := []interface{}{}, map[string][]int{}
		err := tt.item.AppendSQLExclude(tt.dialect, buf, &gotArgs, gotParams, nil, tt.excludedTableQualifiers)
		if err != nil {
			t.Fatal(testutil.Callers(), err)
		}
		if diff := testutil.Diff(buf.String(), tt.wantQuery); diff != "" {
			t.Error(testutil.Callers(), diff)
		}
		if diff := testutil.Diff(gotArgs, tt.wantArgs); diff != "" {
			t.Error(testutil.Callers(), diff)
		}
	}

	t.Run("BooleanField", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.item = NewBooleanField("field", TableInfo{TableName: "tbl"})
		tt.wantQuery = "tbl.field"
		tt.wantArgs = []interface{}{}
		assert(t, tt)
	})

	t.Run("BooleanField with alias", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.item = NewBooleanField("field", TableInfo{TableName: "tbl"}).As("f")
		tt.wantQuery = "tbl.field"
		tt.wantArgs = []interface{}{}
		assert(t, tt)
	})

	t.Run("BooleanField ASC NULLS LAST", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.item = NewBooleanField("field", TableInfo{TableName: "tbl"}).Asc().NullsLast()
		tt.wantQuery = "tbl.field ASC NULLS LAST"
		tt.wantArgs = []interface{}{}
		assert(t, tt)
	})

	t.Run("BooleanField DESC NULLS FIRST", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.item = NewBooleanField("field", TableInfo{TableName: "tbl"}).Desc().NullsFirst()
		tt.wantQuery = "tbl.field DESC NULLS FIRST"
		tt.wantArgs = []interface{}{}
		assert(t, tt)
	})

	t.Run("BooleanField NOT", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.item = NewBooleanField("field", TableInfo{TableName: "tbl"}).Not()
		tt.wantQuery = "NOT tbl.field"
		tt.wantArgs = []interface{}{}
		assert(t, tt)
	})

	t.Run("BooleanField IS NULL", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.item = NewBooleanField("field", TableInfo{TableName: "tbl"}).IsNull()
		tt.wantQuery = "tbl.field IS NULL"
		tt.wantArgs = []interface{}{}
		assert(t, tt)
	})

	t.Run("BooleanField IS NOT NULL", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.item = NewBooleanField("field", TableInfo{TableName: "tbl"}).IsNotNull()
		tt.wantQuery = "tbl.field IS NOT NULL"
		tt.wantArgs = []interface{}{}
		assert(t, tt)
	})

	t.Run("BooleanField Eq", func(t *testing.T) {
		t.Parallel()
		var tt TT
		f := NewBooleanField("field", TableInfo{TableName: "tbl"})
		tt.item = f.Eq(f)
		tt.wantQuery = "tbl.field = tbl.field"
		tt.wantArgs = []interface{}{}
		assert(t, tt)
	})

	t.Run("BooleanField Ne", func(t *testing.T) {
		t.Parallel()
		var tt TT
		f := NewBooleanField("field", TableInfo{TableName: "tbl"})
		tt.item = f.Ne(f)
		tt.wantQuery = "tbl.field <> tbl.field"
		tt.wantArgs = []interface{}{}
		assert(t, tt)
	})

	t.Run("BooleanField EqBool", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.item = NewBooleanField("field", TableInfo{TableName: "tbl"}).EqBool(true)
		tt.wantQuery = "tbl.field = ?"
		tt.wantArgs = []interface{}{true}
		assert(t, tt)
	})

	t.Run("BooleanField NeBool", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.item = NewBooleanField("field", TableInfo{TableName: "tbl"}).NeBool(true)
		tt.wantQuery = "tbl.field <> ?"
		tt.wantArgs = []interface{}{true}
		assert(t, tt)
	})

	t.Run("BooleanField SetBool", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.item = NewBooleanField("field", TableInfo{TableName: "tbl"}).SetBool(true)
		tt.wantQuery = "tbl.field = ?"
		tt.wantArgs = []interface{}{true}
		assert(t, tt)
	})
}
