package eligible

import (
	D "github.com/Netflix/chaosmonkey/deploy"
	"github.com/Netflix/chaosmonkey/grp"
	"github.com/Netflix/chaosmonkey/mock"
	"testing"
	"sort"
	"github.com/Netflix/chaosmonkey"
)

func mockDeployment() D.Deployment {
	a := D.AccountName("prod")
	p := "aws"
	r := D.RegionName("us-east-1")

	return &mock.Deployment{AppMap: map[string]D.AppMap{
		"foo": {a:
		D.AccountInfo{CloudProvider: p, Clusters:
		D.ClusterMap{
			"foo-prod": {r: {"foo-prod-v001": []D.InstanceID{"i-11111111", "i-22222222"}}},
			"foo-prod-lorin": {r: {"foo-prod-lorin-v123": []D.InstanceID{"i-33333333", "i-44444444"}}},
			"foo-staging": {r: {"foo-staging-v005": []D.InstanceID{"i-55555555", "i-66666666"}}},
			"foo-staging-lorin": {r: {"foo-prod-lorin-v117": []D.InstanceID{"i-77777777", "i-88888888"}}},
		}},
		}}}
}

func TestClusterGropuing(t *testing.T) {
	// setup
	dep := mockDeployment()
	group := grp.New("foo", "prod", "us-east-1", "", "foo-prod")

	// code under test
	instances, err := Instances(group, nil, dep)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// assertions
	gots := ids(instances)
	wants := []string{"i-11111111", "i-22222222"}

	if got, want := len(gots), len(wants); got != want {
		t.Fatalf("len(eligible.Instances(group, cfg, app))=%v, want %v", got, want)
	}

	for i, got := range gots {
		if want := wants[i]; got != want {
			t.Fatalf("got=%v, want=%v", got, want)
		}
	}
}

func TestStackGrouping(t *testing.T) {
	// setup
	dep := mockDeployment()
	group := grp.New("foo", "prod", "us-east-1", "staging", "")

	// code under test
	instances, err := Instances(group, nil, dep)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// assertions
	gots := ids(instances)
	wants := []string{"i-55555555", "i-66666666", "i-77777777", "i-88888888"}

	if got, want := len(gots), len(wants); got != want {
		t.Fatalf("len(eligible.Instances(group, cfg, app))=%v, want %v", got, want)
	}

	for i, got := range gots {
		if want := wants[i]; got != want {
			t.Fatalf("got=%v, want=%v", got, want)
		}
	}
}

// ids returns a sorted list of instance ids
func ids(instances []chaosmonkey.Instance) []string {
	result := make([]string, len(instances))
	for i, inst := range instances {
		result[i] = inst.ID()
	}

	sort.Strings(result)
	return result

}

func TestAppGrouping(t *testing.T) {
	// setup
	dep := mockDeployment()
	group := grp.New("foo", "prod", "us-east-1", "", "")

	// code under test
	instances, err := Instances(group, nil, dep)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	// assertions
	gots := ids(instances)
	wants := []string{"i-11111111", "i-22222222", "i-33333333", "i-44444444", "i-55555555", "i-66666666", "i-77777777", "i-88888888"}

	if got, want := len(gots), len(wants); got != want {
		t.Fatalf("len(eligible.Instances(group, cfg, app))=%v, want %v", got, want)
	}

	for i, got := range gots {
		if want := wants[i]; got != want {
			t.Fatalf("got=%v, want=%v", got, want)
		}
	}

}

