package azure

import (
	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	"github.com/Azure/go-autorest/autorest/to"
	"github.com/google/badwolf/triple"
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/rdf"
)

func (inf *Infra) AvailabilitysetsGraph() (*rdf.Graph, error) {
	var triples []*triple.Triple
	gc := resources.NewGroupsClient(AZURE_SUBSCRIPTION_ID)
	gc.Authorizer = inf.Spt
	list, err := gc.List("", nil)
	if err != nil {
		return nil, err
	}
	for _, rg := range *list.Value {
		asc := compute.NewAvailabilitySetsClient(AZURE_SUBSCRIPTION_ID)
		asc.Authorizer = inf.Spt
		r, err := asc.List(to.String(rg.Name))
		if err != nil {
			return nil, err
		}
		for _, as := range *r.Value {
			_, err := cloud.AddNodeWithPropertiesToTriples(AVAILABILITYSET, to.String(as.Name), &as, resourcesProperties, &triples)
			if err != nil {
				return nil, err
			}
		}
	}

	return rdf.NewGraphFromTriples(triples), nil
}

func (inf *Infra) VirtualmachinesGraph() (*rdf.Graph, error) {
	var triples []*triple.Triple
	vmc := compute.NewVirtualMachinesClient(AZURE_SUBSCRIPTION_ID)
	vmc.Authorizer = inf.Spt
	vms, err := vmc.ListAll()
	if err != nil {
		return nil, err
	}
	for _, vm := range *vms.Value {
		vmNode, err := cloud.AddNodeWithPropertiesToTriples(VIRTUALMACHINE, to.String(vm.Name), &vm, resourcesProperties, &triples)
		if err != nil {
			return nil, err
		}
		var publicIP string
		var privateIP string
		if len(*vm.VirtualMachineProperties.NetworkProfile.NetworkInterfaces) > 0 {
			firstNetworkInterface := (*vm.VirtualMachineProperties.NetworkProfile.NetworkInterfaces)[0]
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
		if privateIP != "" {
			if propT, err := cloud.PropertyTriple(vmNode, "PrivateIp", privateIP); err != nil {
				return nil, err
			} else {
				triples = append(triples, propT)
			}
		}
		if publicIP != "" {
			if propT, err := cloud.PropertyTriple(vmNode, "PublicIp", publicIP); err != nil {
				return nil, err
			} else {
				triples = append(triples, propT)
			}
		}
	}

	return rdf.NewGraphFromTriples(triples), nil
}
