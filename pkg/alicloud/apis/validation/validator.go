package validation

import (
	"fmt"
	"strings"

	"github.com/gardener/machine-controller-manager-provider-alicloud/pkg/alicloud/apis"
)

// ValidateProviderSpec validates AlicloudProviderSpec and secrets
func ValidateProviderSpec(spec *api.AlicloudProviderSpec) []error {
	allErrs := []error{}

	if "" == spec.ImageID {
		allErrs = append(allErrs, fmt.Errorf("ImageID is required"))
	}
	if "" == spec.Region {
		allErrs = append(allErrs, fmt.Errorf("Region is required"))
	}
	if "" == spec.ZoneID {
		allErrs = append(allErrs, fmt.Errorf("ZoneID is required"))
	}
	if "" == spec.InstanceType {
		allErrs = append(allErrs, fmt.Errorf("InstanceType is required"))
	}
	if "" == spec.VSwitchID {
		allErrs = append(allErrs, fmt.Errorf("VSwitchID is required"))
	}
	if "" == spec.KeyPairName {
		allErrs = append(allErrs, fmt.Errorf("KeyPairName is required"))
	}

	allErrs = append(allErrs, validateTags(spec.Tags)...)

	return allErrs
}

func validateTags(tags map[string]string) []error {
	allErrs := []error{}
	clusterName := ""
	nodeRole := ""

	if tags == nil {
		allErrs = append(allErrs, fmt.Errorf("Tags required for Alicloud machines"))
	}

	for key := range tags {
		if strings.Contains(key, "kubernetes.io/cluster/") {
			clusterName = key
		} else if strings.Contains(key, "kubernetes.io/role/") {
			nodeRole = key
		}
	}

	if clusterName == "" {
		allErrs = append(allErrs, fmt.Errorf("Tag required of the form kubernetes.io/cluster/****"))
	}
	if nodeRole == "" {
		allErrs = append(allErrs, fmt.Errorf("Tag required of the form kubernetes.io/role/****"))
	}

	return allErrs
}
