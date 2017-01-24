package display

import (
	"bytes"
	"testing"

	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/literal"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/rdf"
)

func TestDisplayCommit(t *testing.T) {
	t0 := parseTriple(`/region<eu-west-1>	"has_type"@[]	"/region"^^type:text`)
	t1 := parseTriple(`/instance<inst_1>	"has_type"@[]	"/instance"^^type:text`)
	t2 := parseTriple(`/instance<inst_1>	"property"@[]	"{"Key":"Id","Value":"inst_1"}"^^type:text`)
	t3 := parseTriple(`/instance<inst_2>	"has_type"@[]	"/instance"^^type:text`)
	t4 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /instance<inst_1>`)
	t5 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /instance<inst_2>`)
	t6 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /vpc<vpc_1>`)
	t7 := parseTriple(`/vpc<vpc_1>  "parent_of"@[] /instance<inst_1>`)
	t7bis := parseTriple(`/instance<inst_1>	"property"@[]	"{"Key":"Deleted","Value":"del_1"}"^^type:text`)
	graph := graph.NewGraphFromTriples([]*triple.Triple{t0, t1, t2, t3, t4, t5, t6, t7, t7bis})
	diff := graph.NewEmptyDiffFromGraph(graph)
	diff.AddDeleted(t2, graph.ParentOfPredicate)
	diff.AddDeleted(t3, graph.ParentOfPredicate)
	diff.AddDeleted(t5, graph.ParentOfPredicate)
	diff.AddDeleted(t7bis, graph.ParentOfPredicate)
	t8 := parseTriple(`/instance<inst_1>	"property"@[]	"{"Key":"Id","Value":"new_id"}"^^type:text`)
	t9 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /instance<inst_3>`)
	t10 := parseTriple(`/instance<inst_3>	"has_type"@[]	"/instance"^^type:text`)
	t11 := parseTriple(`/instance<inst_3>	"property"@[]	"{"Key":"Id","Value":"inst_3"}"^^type:text`)
	t12 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /instance<inst_4>`)
	t13 := parseTriple(`/instance<inst_4>	"has_type"@[]	"/instance"^^type:text`)
	t14 := parseTriple(`/instance<inst_4>	"property"@[]	"{"Key":"Test","Value":"test_1"}"^^type:text`)
	t15 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /instance<inst_5>`)
	t16 := parseTriple(`/instance<inst_5>	"has_type"@[]	"/instance"^^type:text`)
	t17 := parseTriple(`/region<eu-west-1>  "parent_of"@[] /subnet<test_1>`)
	t18 := parseTriple(`/subnet<test_1>	"property"@[]	"{"Key":"prop","Value":"val"}"^^type:text`)
	diff.AddInserted(t8, graph.ParentOfPredicate)
	diff.AddInserted(t9, graph.ParentOfPredicate)
	diff.AddInserted(t10, graph.ParentOfPredicate)
	diff.AddInserted(t11, graph.ParentOfPredicate)
	diff.AddInserted(t12, graph.ParentOfPredicate)
	diff.AddInserted(t13, graph.ParentOfPredicate)
	diff.AddInserted(t14, graph.ParentOfPredicate)
	diff.AddInserted(t15, graph.ParentOfPredicate)
	diff.AddInserted(t16, graph.ParentOfPredicate)
	diff.AddInserted(t17, graph.ParentOfPredicate)
	diff.AddInserted(t18, graph.ParentOfPredicate)

	rootNode, err := node.NewNodeFromStrings(graph.Region.ToRDFString(), "eu-west-1")
	if err != nil {
		t.Fatal(err)
	}
	table, err := tableFromDiff(diff, rootNode, aws.InfraServiceName)
	if err != nil {
		t.Fatal(err)
	}
	var print bytes.Buffer
	table.Fprint(&print)
	expected := `+----------+----------+----------+----------+
|  TYPE ▲  | NAME/ID  | PROPERTY |  VALUE   |
+----------+----------+----------+----------+
| instance | + inst_3 | Id       | + inst_3 |
+          +----------+----------+----------+
|          | + inst_4 | Test     | + test_1 |
+          +----------+----------+----------+
|          | + inst_5 |          |          |
+          +----------+----------+----------+
|          | - inst_2 |          |          |
+          +----------+----------+----------+
|          | inst_1   | Deleted  | - del_1  |
+          +          +----------+----------+
|          |          | Id       | + new_id |
+          +          +          +----------+
|          |          |          | - inst_1 |
+----------+----------+----------+----------+
| subnet   | + test_1 | prop     | + val    |
+----------+----------+----------+----------+
`
	if got, want := print.String(), expected; got != want {
		t.Fatalf("got\n%s\nwant\n%s\n", got, want)
	}
}

func parseTriple(s string) *triple.Triple {
	t, err := triple.Parse(s, literal.DefaultBuilder())
	if err != nil {
		panic(err)
	}

	return t
}
