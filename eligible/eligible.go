// Package eligible contains methods that determine which instances are eligible for Chaos Monkey termination
package eligible

import (
	"github.com/Netflix/chaosmonkey/grp"
	"github.com/Netflix/chaosmonkey"
	"github.com/Netflix/chaosmonkey/deploy"
	"github.com/pkg/errors"
)

// Instances returns instances eligible for termination
func Instances(group grp.InstanceGroup, cfg chaosmonkey.AppConfig, dep deploy.Deployment) ([]*deploy.Instance, error) {
	return nil, errors.New("not yet impelemented")
}