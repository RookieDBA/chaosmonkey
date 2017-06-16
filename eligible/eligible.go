package eligible

import (
	"github.com/Netflix/chaosmonkey/grp"
	"github.com/Netflix/chaosmonkey"
	"github.com/Netflix/chaosmonkey/deploy"
)

// Instances returns instances eligible for termination
func Instances(group grp.InstanceGroup, cfg chaosmonkey.AppConfig, dep deploy.Deployment) []*deploy.Instance {
	return nil
}