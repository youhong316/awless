package azure

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/arm/compute"
	"github.com/Azure/azure-sdk-for-go/arm/network"
	"github.com/Azure/azure-sdk-for-go/arm/resources/resources"
	azuresdk "github.com/Azure/go-autorest/autorest/azure"
	"github.com/Azure/go-autorest/autorest/to"
)

type Infra struct {
	Spt *azuresdk.ServicePrincipalToken
}

var InfraService *Infra

func InitServices(spt *azuresdk.ServicePrincipalToken) {
	InfraService = NewInfra(spt)
}

func (inf *Infra) Instances() (interface{}, error) {
	vmc := compute.NewVirtualMachinesClient(AZURE_SUBSCRIPTION_ID)
	vmc.Authorizer = inf.Spt

	return vmc.ListAll()
}

func (inf *Infra) AvailabilitySets() (interface{}, error) {
	gc := resources.NewGroupsClient(AZURE_SUBSCRIPTION_ID)
	gc.Authorizer = inf.Spt
	result := compute.AvailabilitySetListResult{}
	result.Value = &([]compute.AvailabilitySet{})
	list, err := gc.List("", nil)
	if err != nil {
		return result, err
	}
	for _, rg := range *list.Value {
		asc := compute.NewAvailabilitySetsClient(AZURE_SUBSCRIPTION_ID)
		asc.Authorizer = inf.Spt
		r, err := asc.List(to.String(rg.Name))
		if err != nil {
			return result, err
		}
		*result.Value = append(*result.Value, *r.Value...)
	}

	return result, nil
}

func NewInfra(spt *azuresdk.ServicePrincipalToken) *Infra {
	return &Infra{Spt: spt}
}

func (inf *Infra) FetchNetworkInterfaceInformation(id string) (*network.Interface, error) {
	nic := network.NewInterfacesClient(AZURE_SUBSCRIPTION_ID)
	nic.Authorizer = inf.Spt
	res, err := nic.ListAll()
	if err != nil {
		return nil, err
	}
	for _, n := range *res.Value {
		if to.String(n.ID) == id {
			return &n, nil
		}
	}
	return nil, fmt.Errorf("network interface not found")
}

func (inf *Infra) FetchPublicIP(id string) (*network.PublicIPAddress, error) {
	nic := network.NewPublicIPAddressesClient(AZURE_SUBSCRIPTION_ID)
	nic.Authorizer = inf.Spt
	res, err := nic.ListAll()
	if err != nil {
		return nil, err
	}
	for _, n := range *res.Value {
		if to.String(n.ID) == id {
			return &n, nil
		}
	}
	return nil, fmt.Errorf("network interface not found")
}
