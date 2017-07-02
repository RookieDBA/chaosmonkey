package eligible

import (
	"github.com/Netflix/chaosmonkey"
	"github.com/Netflix/chaosmonkey/grp"
	"github.com/Netflix/chaosmonkey/mock"
	"testing"
)

func TestClusterGropuing(t *testing.T) {
	// setup
	appConfig := chaosmonkey.AppConfig{
		Enabled:                        true,
		RegionsAreIndependent:          true,
		MeanTimeBetweenKillsInWorkDays: 5,
		MinTimeBetweenKillsInWorkDays:  1,
		Grouping:                       chaosmonkey.Cluster,
	}
	dep := mock.Deployment()
	group := grp.New("foo", "prod", "us-east-1", "", "foo-prod")

	// code under test
	instances, err := Instances(group, appConfig.Exceptions, dep)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// assertions
	wants := []string{"i-d3e3d611", "i-63f52e25"}

	if got, want := len(instances), 2; got != want {
		t.Fatalf("len(eligible.Instances(group, cfg, app))=%v, want %v", got, want)
	}

	for i, inst := range instances {
		if got, want := inst.ID(), wants[i]; got != want {
			t.Fatalf("got=%v, want=%v", got, want)
		}
	}
}
