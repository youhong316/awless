package azure

const (
	AVAILABILITYSET = "/availabilityset"
	VIRTUALMACHINE  = "/virtualmachine"
)

var resourcesProperties = map[string]map[string]string{
	AVAILABILITYSET: {
		"Id":       "Name",
		"Name":     "Name",
		"Type":     "Type",
		"Location": "Location",
		"VMs":      "VirtualMachines",
	},
	VIRTUALMACHINE: {
		"Id":    "ID",
		"Name":  "Name",
		"Type":  "Type",
		"State": "ProvisioningState",
	},
}
