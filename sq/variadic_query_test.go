package sq

import (
	"bytes"
	"testing"

	"github.com/bokwoon95/pagemanager/sq/internal/testutil"
)

func TestVariadicQuery(t *testing.T) {
	type TT struct {
		dialect    string
		item       SQLAppender
		wantQuery  string
		wantArgs   []interface{}
		wantParams map[string][]int
	}

	const dialect = DialectMySQL

	assert := func(t *testing.T, tt TT) {
		buf := bufpool.Get().(*bytes.Buffer)
		defer func() {
			buf.Reset()
			bufpool.Put(buf)
		}()
		gotArgs, gotParams := []interface{}{}, map[string][]int{}
		err := tt.item.AppendSQL(tt.dialect, buf, &gotArgs, gotParams, nil)
		if err != nil {
			t.Fatal(testutil.Callers(), err)
		}
		if diff := testutil.Diff(buf.String(), tt.wantQuery); diff != "" {
			t.Error(testutil.Callers(), diff)
		}
		if diff := testutil.Diff(gotArgs, tt.wantArgs); diff != "" {
			t.Error(testutil.Callers(), diff)
		}
		if tt.wantParams != nil {
			if diff := testutil.Diff(gotParams, tt.wantParams); diff != "" {
				t.Error(testutil.Callers(), diff)
			}
		}
	}

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.item = VariadicQuery{}
		tt.wantArgs = []interface{}{}
		assert(t, tt)
	})

	t.Run("nested single variadic queries", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.item = Union(UnionAll(Intersect(IntersectAll(Except(ExceptAll(Queryf("SELECT {}", 1)))))))
		tt.wantQuery = "SELECT ?"
		tt.wantArgs = []interface{}{1}
		assert(t, tt)
	})

	t.Run("nested variadic queries", func(t *testing.T) {
		t.Parallel()
		var tt TT
		tt.item = Intersect(
			Union(Union(Queryf("SELECT {}", 1)), Queryf("SELECT {}", "abc")),
			Union(Queryf("SELECT {}", 3.14)),
		)
		tt.wantQuery = "(SELECT ? UNION SELECT ?) INTERSECT SELECT ?"
		tt.wantArgs = []interface{}{1, "abc", 3.14}
		assert(t, tt)
	})
}
