package alicloud

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	api "github.com/gardener/machine-controller-manager-provider-alicloud/pkg/alicloud/apis"
	validator "github.com/gardener/machine-controller-manager-provider-alicloud/pkg/alicloud/apis/validation"
	"github.com/satori/go.uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func decodeProviderSpecAndSecret(providerSpecBytes []byte,
	secretMap map[string][]byte,
	checkUserData bool) (*api.AlicloudProviderSpec, *api.Secrets, error) {
	var (
		providerSpec api.AlicloudProviderSpec
	)
	// Extract providerSpec
	err := json.Unmarshal(providerSpecBytes, &providerSpec)
	if err != nil {
		return nil, nil, status.Error(codes.Internal, err.Error())
	}

	// Extract secrets from secretMap
	secrets, err := getSecretsFromSecretMap(secretMap, checkUserData)
	if err != nil {
		return nil, nil, err
	}

	if errs := validator.ValidateProviderSpec(&providerSpec); len(errs) > 0 {
		err = fmt.Errorf("Error while validating ProviderSpec %v", errs)
		return nil, nil, status.Error(codes.Internal, err.Error())
	}

	return &providerSpec, secrets, nil
}

// getSecretsFromSecretMap converts secretMap to api.secrets object
func getSecretsFromSecretMap(secretMap map[string][]byte, checkUserData bool) (*api.Secrets, error) {
	accessKeyID, keyIDExists := secretMap["alicloudAccessKeyID"]
	accessKeySecret, accessKeyExists := secretMap["alicloudAccessKeySecret"]
	userData, userDataExists := secretMap["userData"]
	if !keyIDExists || !accessKeyExists || (checkUserData && !userDataExists) {
		var err error
		if checkUserData {
			err = fmt.Errorf(
				"Invalidate Secret Map. Map variables present \nalicloudAccessKeyID: %t, \nalicloudAccessKeySecret: %t, \nuserData: %t",
				keyIDExists,
				accessKeyExists,
				userDataExists,
			)
		} else {
			err = fmt.Errorf(
				"Invalidate Secret Map. Map variables present \nalicloudAccessKeyID: %t, \nalicloudAccessKeySecret: %t",
				keyIDExists,
				accessKeyExists,
			)
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	var secrets api.Secrets
	secrets.AlicloudAccessKeyID = string(accessKeyID)
	secrets.AlicloudAccessKeySecret = string(accessKeySecret)
	secrets.UserData = string(userData)

	return &secrets, nil
}

func encodeProviderID(region, machineID string) string {
	return fmt.Sprintf("%s.%s", region, machineID)
}

func decodeProviderID(id string) (string, error) {
	splitProviderID := strings.Split(id, ".")
	if len(splitProviderID) == 1 {
		return "", fmt.Errorf("Provider ID should be in the format [region].[instance-id], current value is %v", id)
	}
	return splitProviderID[len(splitProviderID)-1], nil
}

func convertToInstanceTags(tags map[string]string) ([]ecs.RunInstancesTag, error) {
	result := []ecs.RunInstancesTag{{}, {}}
	hasCluster := false
	hasRole := false

	for k, v := range tags {
		if strings.Contains(k, "kubernetes.io/cluster/") {
			hasCluster = true
			result[0].Key = k
			result[0].Value = v
		} else if strings.Contains(k, "kubernetes.io/role/") {
			hasRole = true
			result[1].Key = k
			result[1].Value = v
		} else {
			result = append(result, ecs.RunInstancesTag{Key: k, Value: v})
		}
	}

	if !hasCluster || !hasRole {
		return nil, fmt.Errorf("Tags should at least contains 2 keys, which are prefixed with kubernetes.io/cluster and kubernetes.io/role")

	}

	return result, nil
}

func getUUIDV4() string {
	uuidV4 := uuid.NewV4()
	return hex.EncodeToString(uuidV4.Bytes())
}

// Host name in Alicloud has relationship with Instance ID
// i-uf69zddmom11ci7est12 => iZuf69zddmom11ci7est12Z
func idToName(instanceID string) string {
	return strings.Replace(instanceID, "-", "Z", 1) + "Z"
}

func (ms *MachinePlugin) getInstancesFromMachineName(machineName string, providerSpec *api.AlicloudProviderSpec, secrets *api.Secrets) ([]ecs.Instance, error) {
	request, err := buildDescribeInstancesRequestWithTags(providerSpec, secrets)
	if err != nil {
		return nil, err
	}
	request.InstanceName = machineName
	return ms.describeInstances(providerSpec.Region, secrets, request)
}

func (ms *MachinePlugin) describeInstances(region string, secrets *api.Secrets, request *ecs.DescribeInstancesRequest) ([]ecs.Instance, error) {
	client, err := ms.SPI.CreateClient(region, secrets.AlicloudAccessKeyID, secrets.AlicloudAccessKeySecret)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	response, err := client.DescribeInstances(request)
	if err != nil {
		return nil, err
	}

	return response.Instances.Instance, nil
}

func buildDescribeInstancesRequestWithTags(providerSpec *api.AlicloudProviderSpec, secrets *api.Secrets) (*ecs.DescribeInstancesRequest, error) {
	request := ecs.CreateDescribeInstancesRequest()
	searchClusterName := ""
	searchNodeRole := ""
	searchClusterNameValue := ""
	searchNodeRoleValue := ""

	for k, v := range providerSpec.Tags {
		if strings.Contains(k, "kubernetes.io/cluster/") {
			searchClusterName = k
			searchClusterNameValue = v
		} else if strings.Contains(k, "kubernetes.io/role/") {
			searchNodeRole = k
			searchNodeRoleValue = v
		}
	}

	if searchClusterName == "" || searchNodeRole == "" {
		return nil, fmt.Errorf("Can't find VMs with none of machineID/Tag[kubernetes.io/cluster/*]/Tag[kubernetes.io/role/*]")
	}

	request.Tag = &[]ecs.DescribeInstancesTag{
		{Key: searchClusterName, Value: searchClusterNameValue},
		{Key: searchNodeRole, Value: searchNodeRoleValue},
	}

	return request, nil
}

func (ms *MachinePlugin) fetchInstances(providerSpec *api.AlicloudProviderSpec, secrets *api.Secrets) ([]ecs.Instance, error) {
	request, err := buildDescribeInstancesRequestWithTags(providerSpec, secrets)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return ms.describeInstances(providerSpec.Region, secrets, request)
}
