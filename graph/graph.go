package graph

import (
	"github.com/google/badwolf/triple"
	"github.com/google/badwolf/triple/node"
	"github.com/wallix/awless/graph/internal/rdf"
)

type Graph struct {
	*rdf.Graph
}

func NewGraph() *Graph {
	return &Graph{rdf.NewGraph()}
}

func NewGraphFromFile(filepath string) (*Graph, error) {
	g, err := rdf.NewGraphFromFile(filepath)
	return &Graph{g}, err
}

func NewParentOfTriple(subject, obj *node.Node) (*triple.Triple, error) {
	return triple.New(subject, rdf.ParentOfPredicate, triple.NewNodeObject(obj))
}

func NewRegionTypeTriple(subject *node.Node) (*triple.Triple, error) {
	return triple.New(subject, rdf.HasTypePredicate, triple.NewLiteralObject(rdf.RegionLiteral))
}
