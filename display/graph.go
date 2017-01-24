package display

import (
	"fmt"
	"os"

	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/graph"
)

// SeveralResourcesOfGraph prints a RDF graph with different type of resources according to there display properties
func SeveralResourcesOfGraph(graph *graph.Graph, displayer *ServiceDisplayer, onlyIDs bool) {
	table := NewTable([]*PropertyDisplayer{{Property: "Type", DontTruncate: true}, {Property: "Name/Id", DontTruncate: true}, {Property: "Property", DontTruncate: true}, {Property: "Value", DontTruncate: true}})
	table.MergeIdenticalCells = true
	for t := range displayer.Resources {
		nodes, err := graph.NodesForType(t)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		for _, node := range nodes {
			res := cloud.InitFromRdfNode(node)
			err := res.UnmarshalFromGraph(graph)
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			for _, propD := range displayer.Resources[t].Properties {
				table.AddValue("Type", t.String())
				table.AddValue("Name/Id", nameOrID(res))
				table.AddValue("Property", propD.displayName())
				table.AddValue("Value", propD.display(propD.propertyValue(res.Properties())))
			}
		}
	}

	table.SetSortBy("Type", "Name/Id", "Property", "Value")

	if onlyIDs {
		table.FprintColumnValues(os.Stdout, "Name/Id", " ")
	} else {
		table.Fprint(os.Stdout)
	}
}
