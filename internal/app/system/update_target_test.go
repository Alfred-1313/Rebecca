package system

import (
	"context"
	"testing"
)

type fakeBuildCatalog struct {
	catalog BuildCatalog
}

func (f fakeBuildCatalog) Status(context.Context, string, *string, string) UpdateStatus {
	return UpdateStatus{}
}

func (f fakeBuildCatalog) Builds(context.Context, string) (BuildCatalog, error) {
	return f.catalog, nil
}

func TestValidateUpdateTargetRefusesDowngrades(t *testing.T) {
	checker := fakeBuildCatalog{catalog: BuildCatalog{
		Floor:  versionSwitchFloor,
		Stable: []BuildVersion{{Version: "v1.6.0", Channel: "stable"}},
		Dev:    []BuildVersion{{Version: "dev-b6f54d6", Channel: "dev", Commit: "b6f54d6eb6c0ffee1234567890abcdef12345678"}},
	}}
	cases := []struct {
		name    string
		channel string
		version string
		current string
		allowed bool
	}{
		{"release below the floor", "", "v1.5.0", "stable", false},
		{"release at the floor", "", "v1.6.0", "stable", true},
		{"dev build before the schema floor", "", "dev-0123456", "dev", false},
		{"dev build at the schema floor", "", "dev-b6f54d6", "dev", true},
		{"dev channel resolves to a listed dev build", "dev", "", "", true},
		{"stable channel resolves to a listed release", "stable", "", "", true},
		{"current channel follows the running stable channel", "current", "", "stable", true},
		{"current channel with an unknown running channel is refused", "", "", "", false},
		{"unknown release is not installable", "", "v9.9.9", "stable", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateUpdateTarget(context.Background(), checker, PanelUpdateRepo, tc.channel, tc.version, tc.current)
			if (err == nil) != tc.allowed {
				t.Fatalf("ValidateUpdateTarget(channel=%q, version=%q, current=%q) error = %v, want allowed=%v", tc.channel, tc.version, tc.current, err, tc.allowed)
			}
		})
	}
}

func TestValidateUpdateTargetRefusesWhenNothingIsAboveTheFloor(t *testing.T) {
	checker := fakeBuildCatalog{catalog: BuildCatalog{Floor: versionSwitchFloor}}
	for _, tc := range []struct{ channel, version string }{
		{"stable", ""},
		{"dev", ""},
		{"", "dev-b6f54d6"},
	} {
		if err := ValidateUpdateTarget(context.Background(), checker, PanelUpdateRepo, tc.channel, tc.version, "stable"); err == nil {
			t.Fatalf("empty catalog must refuse channel=%q version=%q", tc.channel, tc.version)
		}
	}
}

func TestNodeUpdatesUseTheNodeFloor(t *testing.T) {
	if floorForRepo(NodeUpdateRepo) != nodeVersionFloor || floorForRepo(PanelUpdateRepo) != versionSwitchFloor {
		t.Fatalf("floorForRepo must pick the floor of each repository")
	}
	checker := fakeBuildCatalog{catalog: BuildCatalog{
		Floor:  nodeVersionFloor,
		Stable: []BuildVersion{{Version: "v1.6.0", Channel: "stable"}},
	}}
	if err := ValidateUpdateTarget(context.Background(), checker, NodeUpdateRepo, "", "v1.5.0", "stable"); err == nil {
		t.Fatal("node target below the node floor must be refused")
	}
	if err := ValidateUpdateTarget(context.Background(), checker, NodeUpdateRepo, "", "v1.6.0", "stable"); err != nil {
		t.Fatalf("node target at the node floor must be allowed: %v", err)
	}
}

func TestUpdateFollowsCurrentChannel(t *testing.T) {
	if !UpdateFollowsCurrentChannel("", "") || !UpdateFollowsCurrentChannel("current", "") {
		t.Fatal("an empty or current channel request must follow the running channel")
	}
	if UpdateFollowsCurrentChannel("", "v1.6.0") || UpdateFollowsCurrentChannel("stable", "") || UpdateFollowsCurrentChannel("dev", "") {
		t.Fatal("an explicit build or channel must not follow the running channel")
	}
}
