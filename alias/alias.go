package alias

import "github.com/wallix/awless/cloud"

type Alias string

func (a Alias) ResolveToId(g *graph.Graph, resT graph.ResourceType) (string, bool) {
	prop := cloud.Property{Key: "Name", Value: a}
	propL, err := prop.ToLiteralObject()
	if err != nil {
		return "", false
	}
	triples, err := g.TriplesForPredicateObject(graph.PropertyPredicate, propL)
	if err != nil {
		return "", false
	}
	for _, t := range triples {
		s := t.Subject()
		if s.Type().String() == resT.ToRDFString() {
			return s.ID().String(), true
		}
	}

	return "", false
}
