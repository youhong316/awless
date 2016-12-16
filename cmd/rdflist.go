package cmd

import (
	"fmt"
	"reflect"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/aws"
	"github.com/wallix/awless/cloud/azure"
	"github.com/wallix/awless/rdf"
)

var infraResourcesToDisplay = map[string][]string{
	"instance": []string{"Id", "Tags[].Name", "State.Name", "Type", "PublicIp", "PrivateIp"},
	"vpc":      []string{"Id", "IsDefault", "State", "CidrBlock"},
	"subnet":   []string{"Id", "MapPublicIpOnLaunch", "State", "CidrBlock"},
}

var azureInfraResourcesToDisplay = map[string][]string{
	"availabilityset": []string{"Id", "Name", "Type", "Location", "VMs[]length"},
	"virtualmachine":  []string{"Id", "Name", "State", "Type", "PrivateIp", "PublicIp"},
}

func init() {
	RootCmd.AddCommand(rdfListCmd)
	for resource, properties := range infraResourcesToDisplay {
		rdfListCmd.AddCommand(rdfListInfraResourceCmd(resource, properties))
	}
	for resource, properties := range azureInfraResourcesToDisplay {
		rdfListCmd.AddCommand(rdfListAzureInfraResourceCmd(resource, properties))
	}
}

var rdfListCmd = &cobra.Command{
	Use:   "rdflist",
	Short: "List various type of items: instances, vpc, subnet ...",
}

var rdfListInfraResourceCmd = func(resource string, properties []string) *cobra.Command {
	resources := resource + "s"
	nodeType := "/" + resource
	return &cobra.Command{
		Use:   resources,
		Short: "List " + resources,

		Run: func(cmd *cobra.Command, args []string) {
			fnName := fmt.Sprintf("%sGraph", humanize(resources))
			method := reflect.ValueOf(aws.InfraService).MethodByName(fnName)
			if method.IsValid() && !method.IsNil() {
				methodI := method.Interface()
				if graphFn, ok := methodI.(func() (*rdf.Graph, error)); ok {
					graph, err := graphFn()
					displayGraph(graph, nodeType, properties, err)
					return
				}
			}
			fmt.Println(fmt.Errorf("Unknown type of resource: %s", resource))
			return
		},
	}
}

var rdfListAzureInfraResourceCmd = func(resource string, properties []string) *cobra.Command {
	resources := resource + "s"
	nodeType := "/" + resource
	return &cobra.Command{
		Use:   "azure-" + resources,
		Short: "List Azure " + resources,

		Run: func(cmd *cobra.Command, args []string) {
			fnName := fmt.Sprintf("%sGraph", humanize(resources))
			method := reflect.ValueOf(azure.InfraService).MethodByName(fnName)
			if method.IsValid() && !method.IsNil() {
				methodI := method.Interface()
				if graphFn, ok := methodI.(func() (*rdf.Graph, error)); ok {
					graph, err := graphFn()
					displayGraph(graph, nodeType, properties, err)
					return
				}
			}
			fmt.Println(fmt.Errorf("Unknown type of resource: %s", resource))
			return
		},
	}
}
