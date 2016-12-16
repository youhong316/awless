package azure

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/arm/examples/helpers"
	azuresdk "github.com/Azure/go-autorest/autorest/azure"
)

var (
	AZURE_CLIENT_ID       = os.Getenv("AZURE_CLIENT_ID")
	AZURE_CLIENT_SECRET   = os.Getenv("AZURE_CLIENT_SECRET")
	AZURE_SUBSCRIPTION_ID = os.Getenv("AZURE_SUBSCRIPTION_ID")
	AZURE_TENANT_ID       = os.Getenv("AZURE_TENANT_ID")
)

func InitSession() (*azuresdk.ServicePrincipalToken, error) {
	azureEnvs := map[string]string{
		"AZURE_CLIENT_ID":       AZURE_CLIENT_ID,
		"AZURE_CLIENT_SECRET":   AZURE_CLIENT_SECRET,
		"AZURE_SUBSCRIPTION_ID": AZURE_SUBSCRIPTION_ID,
		"AZURE_TENANT_ID":       AZURE_TENANT_ID,
	}
	err := checkEnvVar(&azureEnvs)
	if err != nil {
		return nil, err
	}
	return helpers.NewServicePrincipalTokenFromCredentials(azureEnvs, azuresdk.PublicCloud.ResourceManagerEndpoint)

}

func checkEnvVar(envVars *map[string]string) error {
	var missingVars []string
	for varName, value := range *envVars {
		if value == "" {
			missingVars = append(missingVars, varName)
		}
	}
	if len(missingVars) > 0 {
		return fmt.Errorf("Missing environment variables %v", missingVars)
	}
	return nil
}
