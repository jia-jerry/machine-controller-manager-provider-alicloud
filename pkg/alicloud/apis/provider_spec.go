/*
Copyright (c) 2019 SAP SE or an SAP affiliate company. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package api

// AlicloudProviderSpec contains the fields of provider spec that the plugin expects
type AlicloudProviderSpec struct {
	// APIVersion mentions the APIVersion of the object being passed
	APIVersion string

	ImageID                 string              `json:"imageID"`
	InstanceType            string              `json:"instanceType"`
	Region                  string              `json:"region"`
	ZoneID                  string              `json:"zoneID,omitempty"`
	SecurityGroupID         string              `json:"securityGroupID,omitempty"`
	VSwitchID               string              `json:"vSwitchID"`
	PrivateIPAddress        string              `json:"privateIPAddress,omitempty"`
	SystemDisk              *AlicloudSystemDisk `json:"systemDisk,omitempty"`
	InstanceChargeType      string              `json:"instanceChargeType,omitempty"`
	InternetChargeType      string              `json:"internetChargeType,omitempty"`
	InternetMaxBandwidthIn  *int                `json:"internetMaxBandwidthIn,omitempty"`
	InternetMaxBandwidthOut *int                `json:"internetMaxBandwidthOut,omitempty"`
	SpotStrategy            string              `json:"spotStrategy,omitempty"`
	IoOptimized             string              `json:"IoOptimized,omitempty"`
	Tags                    map[string]string   `json:"tags,omitempty"`
	KeyPairName             string              `json:"keyPairName"`
}

// AlicloudSystemDisk describes SystemDisk for Alicloud.
type AlicloudSystemDisk struct {
	Category string `json:"category"`
	Size     int    `json:"size"`
}

// Secrets stores the cloud-provider specific sensitive-information.
type Secrets struct {
	// AliCloud access key id (base64 encoded)
	AlicloudAccessKeyID string `json:"alicloudAccessKeyID,omitempty"`
	// AliCloud access key secret(base64 encoded)
	AlicloudAccessKeySecret string `json:"alicloudAccessKeySecret,omitempty"`
	// AliCloud cloud config file (base64 encoded)
	UserData string `json:"userData,omitempty"`
}
