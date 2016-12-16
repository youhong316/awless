package azure

import (
	"fmt"
	"io"
	"path"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/go-autorest/autorest/to"
)

func TabularDisplay(item interface{}, w io.Writer) {
	switch ii := item.(type) {
	case compute.VirtualMachineListResult:
		fmt.Fprintln(w, "Id\tName\tState\tType\tPriv IP\tPub IP")
		for _, inst := range *ii.Value {
			var publicIP string
			var privateIP string
			if len(*inst.VirtualMachineProperties.NetworkProfile.NetworkInterfaces) > 0 {
				firstNetworkInterface := (*inst.VirtualMachineProperties.NetworkProfile.NetworkInterfaces)[0]
				netIntInfo, err := InfraService.FetchNetworkInterfaceInformation(to.String(firstNetworkInterface.ID))
				if err == nil {
					configs := netIntInfo.IPConfigurations
					for _, c := range *configs {
						if to.Bool(c.Primary) {
							privateIP = to.String(c.PrivateIPAddress)
							ip, err := InfraService.FetchPublicIP(to.String(c.PublicIPAddress.ID))
							if err == nil {
								publicIP = to.String(ip.IPAddress)
							}
						}
					}
				}
			}
			_, cleanId := path.Split(to.String(inst.ID))
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s", cleanId, to.String(inst.Name), to.String(inst.VirtualMachineProperties.ProvisioningState), inst.HardwareProfile.VMSize, privateIP, publicIP))
		}
	case compute.AvailabilitySetListResult:
		fmt.Fprintln(w, "Id\tName\tType\tLocation\tNb VMs")

		for _, as := range *ii.Value {
			_, cleanId := path.Split(to.String(as.ID))
			fmt.Fprintln(w, fmt.Sprintf("%s\t%s\t%s\t%s\t%d", cleanId, to.String(as.Name), to.String(as.Type), to.String(as.Location), len(*as.VirtualMachines)))
		}
	default:
		fmt.Printf("%T -> %v\n", item, item)
		return
	}
}
