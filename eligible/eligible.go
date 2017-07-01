// Package eligible contains methods that determine which instances are eligible for Chaos Monkey termination
package eligible

import (
	"github.com/Netflix/chaosmonkey"
	"github.com/Netflix/chaosmonkey/deploy"
	"github.com/Netflix/chaosmonkey/grp"
	"github.com/SmartThingsOSS/frigga-go"
	"github.com/pkg/errors"
)

type (
	cluster struct {
		appName       deploy.AppName
		accountName   deploy.AccountName
		cloudProvider deploy.CloudProvider
		regionName    deploy.RegionName
		clusterName   deploy.ClusterName
	}

	instance struct {
		appName       deploy.AppName
		accountName   deploy.AccountName
		regionName    deploy.RegionName
		stackName     deploy.StackName
		clusterName   deploy.ClusterName
		asgName       deploy.ASGName
		id            deploy.InstanceID
		cloudProvider deploy.CloudProvider
	}
)

func (i instance) AppName() string {
	return string(i.appName)
}

func (i instance) AccountName() string {
	return string(i.accountName)
}

func (i instance) RegionName() string {
	return string(i.regionName)
}

func (i instance) StackName() string {
	return string(i.stackName)
}

func (i instance) ClusterName() string {
	return string(i.clusterName)
}

func (i instance) ASGName() string {
	return string(i.asgName)
}

func (i instance) Name() string {
	return string(i.clusterName)
}

func (i instance) ID() string {
	return string(i.id)
}

func (i instance) CloudProvider() string {
	return string(i.cloudProvider)
}

func isException(exs []chaosmonkey.Exception, account deploy.AccountName, names *frigga.Names, region deploy.RegionName) bool {
	for _, ex := range exs {
		if ex.Matches(string(account), names.Stack, names.Detail, string(region)) {
			return true
		}
	}

	return false
}

func clusters(group grp.InstanceGroup, cloudProvider deploy.CloudProvider, exs []chaosmonkey.Exception, dep deploy.Deployment) ([]cluster, error) {
	account := deploy.AccountName(group.Account())
	clusterNames, err := dep.GetClusterNames(group.App(), account)
	regions := make([]deploy.RegionName, 0)
	region, ok := group.Region()
	if ok {
		regions = append(regions, deploy.RegionName(region))
	} else {
		regions = append(regions, "us-east-1", "us-west-2", "eu-west-1")
	}
	if err != nil {
		return nil, err
	}

	result := make([]cluster, 0)
	for _, clusterName := range clusterNames {
		names, err := frigga.Parse(string(clusterName))
		if err != nil {
			return nil, err
		}

		for _, region := range regions {

			if isException(exs, account, names, region) {
				continue
			}

			if grp.Kontains(group, string(account), string(region), string(clusterName)) {
				result = append(result, cluster{appName: deploy.AppName(names.App),
					accountName:   account,
					cloudProvider: cloudProvider,
					regionName:    region,
					clusterName:   clusterName,
				})
			}
		}
	}

	return result, nil

}

// Instances returns instances eligible for termination
func Instances(group grp.InstanceGroup, cfg chaosmonkey.AppConfig, dep deploy.Deployment) ([]chaosmonkey.Instance, error) {
	if !cfg.Enabled {
		return nil, nil
	}

	if cfg.Whitelist != nil {
		return nil, errors.New("whitelist is not supported")
	}

	cloudProvider, err := dep.CloudProvider(group.Account())
	if err != nil {
		return nil, errors.Wrap(err, "retrieve cloud provider failed")
	}

	cls, err := clusters(group, deploy.CloudProvider(cloudProvider), cfg.Exceptions, dep)
	if err != nil {
		return nil, err
	}

	result := make([]chaosmonkey.Instance, 0)

	for _, cl := range cls {
		instances, err := getInstances(cl, dep)
		if err != nil {
			return nil, err
		}
		result = append(result, instances...)

	}
	return result, nil

}

func getInstances(cl cluster, dep deploy.Deployment) ([]chaosmonkey.Instance, error) {
	result := make([]chaosmonkey.Instance, 0)

	asgName, ids, err := dep.GetInstanceIDs(string(cl.appName), cl.accountName, string(cl.cloudProvider), cl.regionName, cl.clusterName)

	if err != nil {
		return nil, err
	}

	for _, id := range ids {
		names, err := frigga.Parse(string(asgName))
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse")
		}
		result = append(result,
			instance{appName: cl.appName,
				accountName:   cl.accountName,
				regionName:    cl.regionName,
				stackName:     deploy.StackName(names.Stack),
				clusterName:   cl.clusterName,
				asgName:       deploy.ASGName(asgName),
				id:            id,
				cloudProvider: cl.cloudProvider,
			})
	}

	return result, nil
}
