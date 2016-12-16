package cmd

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/wallix/awless/cloud/azure"
)

func init() {
	listAzureCmd.AddCommand(listAzureInstancesCmd)
	listAzureCmd.AddCommand(listAzureAvailabilitySetsCmd)

	RootCmd.AddCommand(listAzureCmd)
}

var listAzureCmd = &cobra.Command{
	Use:   "list-azure",
	Short: "List various type of items for Azure Cloud: users, groups, instances, ...",
}

var listAzureInstancesCmd = &cobra.Command{
	Use:   "instances",
	Short: "List azure instances",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := azure.InfraService.Instances()
		azureDisplay(resp, err, displayFormat)
	},
}

var listAzureAvailabilitySetsCmd = &cobra.Command{
	Use:   "availability-sets",
	Short: "List azure availability sets",

	Run: func(cmd *cobra.Command, args []string) {
		resp, err := azure.InfraService.AvailabilitySets()
		azureDisplay(resp, err, displayFormat)
	},
}

func azureDisplay(item interface{}, err error, format ...string) {
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if len(format) > 0 {
		switch format[0] {
		case "raw":
			fmt.Println(item)
		default:
			azureLineDisplay(item)
		}
	} else {
		azureLineDisplay(item)
	}
}

func azureLineDisplay(item interface{}) {
	w := tabwriter.NewWriter(os.Stdout, 20, 1, 1, ' ', 0)
	azure.TabularDisplay(item, w)
	w.Flush()
}
