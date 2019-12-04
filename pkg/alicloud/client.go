package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

//Client is a client to operater Machine and related resources from AliCloud
type Client interface {
	RunInstances(request *ecs.RunInstancesRequest) (response *ecs.RunInstancesResponse, err error)
	DescribeInstances(request *ecs.DescribeInstancesRequest) (response *ecs.DescribeInstancesResponse, err error)
	DeleteInstance(request *ecs.DeleteInstanceRequest) (response *ecs.DeleteInstanceResponse, err error)
	StopInstance(request *ecs.StopInstanceRequest) (response *ecs.StopInstanceResponse, err error)
}

type clientImpl struct {
	*ecs.Client
}
