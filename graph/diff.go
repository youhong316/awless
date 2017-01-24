package graph

import "github.com/wallix/awless/graph/internal/rdf"

type Diff struct {
	*rdf.Diff
}

var DefaultDiffer = rdf.DefaultDiffer
