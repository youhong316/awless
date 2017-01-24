package cloud

import "github.com/wallix/awless/rdf"

type Service interface {
	FetchRDFResources(graph.ResourceType) (*rdf.Graph, error)
}
