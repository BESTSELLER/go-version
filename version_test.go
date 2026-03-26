// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package version

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestNewVersion(t *testing.T) {
	cases := []struct {
		version string
		err     bool
	}{
		{"", true},
		{"1.2.3", false},
		{"1.0", false},
		{"1", false},
		{"1.2.beta", true},
		{"1.21.beta", true},
		{"foo", true},
		{"1.2-5", false},
		{"1.2-beta.5", false},
		{"\n1.2", true},
		{"1.2.0-x.Y.0+metadata", false},
		{"1.2.0-x.Y.0+metadata-width-hyphen", false},
		{"1.2.3-rc1-with-hyphen", false},
		{"1.2.3.4", false},
		{"1.2.0.4-x.Y.0+metadata", false},
		{"1.2.0.4-x.Y.0+metadata-width-hyphen", false},
		{"1.2.0-X-1.2.0+metadata~dist", false},
		{"1.2.3.4-rc1-with-hyphen", false},
		{"1.2.3.4", false},
		{"v1.2.3", false},
		{"foo1.2.3", true},
		{"1.7rc2", false},
		{"v1.7rc2", false},
		{"1.0-", false},
		{"controller-v0.40.2", true},
		{"azure-cli-v1.4.2", true},
	}

	for _, tc := range cases {
		_, err := NewVersion(tc.version)
		if tc.err && err == nil {
			t.Fatalf("expected error for version: %q", tc.version)
		} else if !tc.err && err != nil {
			t.Fatalf("error for version %q: %s", tc.version, err)
		}
	}
}

func TestNewVersionWithPrefix(t *testing.T) {
	cases := []struct {
		version string
		prefix  string
		err     bool
	}{
		{"", "release-", true},
		{"rel-1.2.3", "release-", true},
		{"release_1.2.3", "release-", true},
		{"release_1.2.0-x.Y.0+metadata", "release_", false},
		{"release-1.2.0-x.Y.0+metadata-width-hyphen", "release-", false},
		{"myrelease-1.2.3-rc1-with-hyphen", "myrelease-", false},
		{"prefix-1.2.3.4", "prefix-", false},
		{"controller-v0.40.2", "controller-", false},
		{"azure-cli-v1.4.2", "azure-cli-", false},
	}

	for _, tc := range cases {
		_, err := NewVersion(tc.version, WithPrefix(tc.prefix))
		if tc.err && err == nil {
			t.Fatalf("expected error for version: %q", tc.version)
		} else if !tc.err && err != nil {
			t.Fatalf("error for version %q: %s", tc.version, err)
		}
	}
}

func TestNewSemver(t *testing.T) {
	cases := []struct {
		version string
		err     bool
	}{
		{"", true},
		{"1.2.3", false},
		{"1.0", false},
		{"1", false},
		{"1.2.beta", true},
		{"1.21.beta", true},
		{"foo", true},
		{"1.2-5", false},
		{"1.2-beta.5", false},
		{"\n1.2", true},
		{"1.2.0-x.Y.0+metadata", false},
		{"1.2.0-x.Y.0+metadata-width-hyphen", false},
		{"1.2.3-rc1-with-hyphen", false},
		{"1.2.3.4", false},
		{"1.2.0.4-x.Y.0+metadata", false},
		{"1.2.0.4-x.Y.0+metadata-width-hyphen", false},
		{"1.2.0-X-1.2.0+metadata~dist", false},
		{"1.2.3.4-rc1-with-hyphen", false},
		{"1.2.3.4", false},
		{"v1.2.3", false},
		{"foo1.2.3", true},
		{"1.7rc2", true},
		{"v1.7rc2", true},
		{"1.0-", true},
		{"controller-v0.40.2", true},
		{"azure-cli-v1.4.2", true},
	}

	for _, tc := range cases {
		_, err := NewSemver(tc.version)
		if tc.err && err == nil {
			t.Fatalf("expected error for version: %q", tc.version)
		} else if !tc.err && err != nil {
			t.Fatalf("error for version %q: %s", tc.version, err)
		}
	}
}

func TestCore(t *testing.T) {
	cases := []struct {
		v1 string
		v2 string
	}{
		{"1.2.3", "1.2.3"},
		{"2.3.4-alpha1", "2.3.4"},
		{"3.4.5alpha1", "3.4.5"},
		{"1.2.3-2", "1.2.3"},
		{"4.5.6-beta1+meta", "4.5.6"},
		{"5.6.7.1.2.3", "5.6.7"},
	}

	for _, tc := range cases {
		v1, err := NewVersion(tc.v1)
		if err != nil {
			t.Fatalf("error for version %q: %s", tc.v1, err)
		}
		v2, err := NewVersion(tc.v2)
		if err != nil {
			t.Fatalf("error for version %q: %s", tc.v2, err)
		}

		actual := v1.Core()
		expected := v2

		if !reflect.DeepEqual(actual, expected) {
			t.Fatalf("expected: %s\nactual: %s", expected, actual)
		}
	}
}

func TestVersionCompare(t *testing.T) {
	cases := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.2.3", "1.4.5", -1},
		{"1.2-beta", "1.2-beta", 0},
		{"1.2", "1.1.4", 1},
		{"1.2", "1.2-beta", 1},
		{"1.2+foo", "1.2+beta", 0},
		{"v1.2", "v1.2-beta", 1},
		{"v1.2+foo", "v1.2+beta", 0},
		{"v1.2.3.4", "v1.2.3.4", 0},
		{"v1.2.0.0", "v1.2", 0},
		{"v1.2.0.0.1", "v1.2", 1},
		{"v1.2", "v1.2.0.0", 0},
		{"v1.2", "v1.2.0.0.1", -1},
		{"v1.2.0.0", "v1.2.0.0.1", -1},
		{"v1.2.3.0", "v1.2.3.4", -1},
		{"1.7rc2", "1.7rc1", 1},
		{"1.7rc2", "1.7", -1},
		{"1.2.0", "1.2.0-X-1.2.0+metadata~dist", 1},
	}

	for _, tc := range cases {
		v1, err := NewVersion(tc.v1)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		v2, err := NewVersion(tc.v2)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v1.Compare(v2)
		expected := tc.expected
		if actual != expected {
			t.Fatalf(
				"%s <=> %s\nexpected: %d\nactual: %d",
				tc.v1, tc.v2,
				expected, actual)
		}
	}
}

func TestVersionCompareWithPrefix(t *testing.T) {
	cases := []struct {
		v1       string
		v1Prefix string
		v2       string
		v2Prefix string
		expected int
	}{
		{"controller-v0.40.2", "controller-", "controller-v0.40.3", "controller-", -1},
		{"0.40.4", "", "controller-v0.40.2", "controller-", 1},
		{"0.40.4", "", "controller-v0.40.4", "controller-", 0},
		{"azure-cli-v1.4.2", "azure-cli-", "azure-cli-v1.4.2", "azure-cli-", 0},
		{"azure-cli-v1.4.1", "azure-cli-", "azure-cli-v1.4.2", "azure-cli-", -1},
		{"1.4.3", "", "azure-cli-v1.4.2", "azure-cli-", 1},
		{"v1.4.3", "", "azure-cli-v1.4.2", "azure-cli-", 1},
		{"controller-v1.4.1", "controller-", "azure-cli-v1.4.2", "azure-cli-", -1},
	}

	for _, tc := range cases {
		var v1 *Version
		var err error
		if tc.v1Prefix != "" {
			v1, err = NewVersion(tc.v1, WithPrefix(tc.v1Prefix))
		} else {
			v1, err = NewVersion(tc.v1)
		}
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		var v2 *Version
		if tc.v2Prefix != "" {
			v2, err = NewVersion(tc.v2, WithPrefix(tc.v2Prefix))
		} else {
			v2, err = NewVersion(tc.v2)
		}
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v1.Compare(v2)
		expected := tc.expected
		if actual != expected {
			t.Fatalf(
				"%s <=> %s\nexpected: %d\nactual: %d",
				tc.v1, tc.v2,
				expected, actual)
		}
	}
}

func TestVersionAccessorsWithPrefix(t *testing.T) {
	v, err := NewVersion("controller-v1.2.0-beta.2+build.5", WithPrefix("controller-"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	if got := v.Prefix(); got != "controller-" {
		t.Fatalf("expected prefix %q, got %q", "controller-", got)
	}

	if got := v.Original(); got != "controller-v1.2.0-beta.2+build.5" {
		t.Fatalf("expected original %q, got %q", "controller-v1.2.0-beta.2+build.5", got)
	}

	if got := v.String(); got != "1.2.0-beta.2+build.5" {
		t.Fatalf("expected string %q, got %q", "1.2.0-beta.2+build.5", got)
	}

	if got := v.Metadata(); got != "build.5" {
		t.Fatalf("expected metadata %q, got %q", "build.5", got)
	}

	if got := v.Prerelease(); got != "beta.2" {
		t.Fatalf("expected prerelease %q, got %q", "beta.2", got)
	}

	expectedSegments := []int{1, 2, 0}
	if got := v.Segments(); !reflect.DeepEqual(got, expectedSegments) {
		t.Fatalf("expected segments %#v, got %#v", expectedSegments, got)
	}

	expectedSegments64 := []int64{1, 2, 0}
	if got := v.Segments64(); !reflect.DeepEqual(got, expectedSegments64) {
		t.Fatalf("expected segments64 %#v, got %#v", expectedSegments64, got)
	}
}

func TestVersionSegmentsWithPrefix(t *testing.T) {
	v, err := NewVersion("azure-cli-v1.4.2", WithPrefix("azure-cli-"))
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	expected := []int{1, 4, 2}
	actual := v.Segments()
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected: %#v\nactual: %#v", expected, actual)
	}
}

func TestVersionCompare_versionAndSemver(t *testing.T) {
	cases := []struct {
		versionRaw string
		semverRaw  string
		expected   int
	}{
		{"0.0.2", "0.0.2", 0},
		{"1.0.2alpha", "1.0.2-alpha", 0},
		{"v1.2+foo", "v1.2+beta", 0},
		{"v1.2", "v1.2+meta", 0},
		{"1.2", "1.2-beta", 1},
		{"v1.2", "v1.2-beta", 1},
		{"1.2.3", "1.4.5", -1},
		{"v1.2", "v1.2.0.0.1", -1},
		{"v1.0.3-", "v1.0.3", -1},
	}

	for _, tc := range cases {
		ver, err := NewVersion(tc.versionRaw)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		semver, err := NewSemver(tc.semverRaw)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := ver.Compare(semver)
		if actual != tc.expected {
			t.Fatalf(
				"%s <=> %s\nexpected: %d\n actual: %d",
				tc.versionRaw, tc.semverRaw, tc.expected, actual,
			)
		}
	}
}

func TestVersionEqual_nil(t *testing.T) {
	mustVersion := func(v string) *Version {
		ver, err := NewVersion(v)
		if err != nil {
			t.Fatal(err)
		}
		return ver
	}
	cases := []struct {
		leftVersion  *Version
		rightVersion *Version
		expected     bool
	}{
		{mustVersion("1.0.0"), nil, false},
		{nil, mustVersion("1.0.0"), false},
		{nil, nil, true},
	}

	for _, tc := range cases {
		given := tc.leftVersion.Equal(tc.rightVersion)
		if given != tc.expected {
			t.Fatalf("expected Equal to nil to be %t", tc.expected)
		}
	}
}

func TestComparePreReleases(t *testing.T) {
	cases := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.2-beta.2", "1.2-beta.2", 0},
		{"1.2-beta.1", "1.2-beta.2", -1},
		{"1.2-beta.2", "1.2-beta.11", -1},
		{"3.2-alpha.1", "3.2-alpha", 1},
		{"1.2-beta.2", "1.2-beta.1", 1},
		{"1.2-beta.11", "1.2-beta.2", 1},
		{"1.2-beta", "1.2-beta.3", -1},
		{"1.2-alpha", "1.2-beta.3", -1},
		{"1.2-beta", "1.2-alpha.3", 1},
		{"3.0-alpha.3", "3.0-rc.1", -1},
		{"3.0-alpha3", "3.0-rc1", -1},
		{"3.0-alpha.1", "3.0-alpha.beta", -1},
		{"5.4-alpha", "5.4-alpha.beta", 1},
		{"v1.2-beta.2", "v1.2-beta.2", 0},
		{"v1.2-beta.1", "v1.2-beta.2", -1},
		{"v3.2-alpha.1", "v3.2-alpha", 1},
		{"v3.2-rc.1-1-g123", "v3.2-rc.2", 1},
	}

	for _, tc := range cases {
		v1, err := NewVersion(tc.v1)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		v2, err := NewVersion(tc.v2)
		if err != nil {
			t.Fatalf("err: %s", err)
		}

		actual := v1.Compare(v2)
		expected := tc.expected
		if actual != expected {
			t.Fatalf(
				"%s <=> %s\nexpected: %d\nactual: %d",
				tc.v1, tc.v2,
				expected, actual)
		}
	}
}

func TestVersionMetadata(t *testing.T) {
	cases := []struct {
		version  string
		expected string