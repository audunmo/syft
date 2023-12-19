package cpp

import (
	"encoding/json"
	"reflect"
	"testing"
)

// The point of these tests is to make sure that we correctly handle the face that `license` is either null or string. We test some other stuff, but that's the important bit
func TestVcpkgFeatureUnmarshal(t *testing.T) {
	tests := []struct {
		name       string
		in         json.RawMessage
		expFeature feature
	}{
		{
			name:       "no license",
			in:         []byte(`{"description": "foo"}`),
			expFeature: feature{Description: "foo"},
		},
		{
			name:       "with license string",
			in:         []byte(`{"description": "foo", "license": "fake-license-v13.36-rc2"}`),
			expFeature: feature{Description: "foo", License: &[]string{"fake-license-v13.36-rc2"}[0]},
		},
		{
			name:       "with null license",
			in:         []byte(`{"description": "foo", "license": null}`),
			expFeature: feature{Description: "foo", License: nil},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var f feature
			err := json.Unmarshal(tc.in, &f)
			if err != nil {
				t.Fatal(err)
			}

			if f.Description != tc.expFeature.Description {
				t.Fatalf("expected description to be %s, got %s", tc.expFeature.Description, f.Description)
			}

			// compare the license values if the pointers are not nil. Both should either be nil, or both should be equal in value
			if f.License != nil && tc.expFeature.License != nil {
				if *f.License != *tc.expFeature.License {
					t.Fatalf("expected license to be %v, got %v", *tc.expFeature.License, *f.License)
				}
				// explicitly test the case where we _want_ the licnse to be nil. If one is nil and the other not, something is wrong, so fail the test
			} else if f.License != nil || tc.expFeature.License != nil {
				t.Fatalf("expected license to be %v, got %v", tc.expFeature.License, f.License)
			}

			if !reflect.DeepEqual(f.Dependencies, tc.expFeature.Dependencies) {
				t.Fatalf("expected dependecies to be %+#v, got %+#v", f.Dependencies, tc.expFeature.Dependencies)
			}
		})
	}
}

func TestVcpkgDependencyUnmarshal(t *testing.T) {
	tests := []struct {
		name   string
		in     json.RawMessage
		expDep vcpkgDependency
	}{
		{
			name:   "json encoded dependency",
			in:     []byte(`{"name": "foo"}`),
			expDep: vcpkgDependency{Name: "foo"},
		},
		{
			name:   "json encoded dependency with MinVersion",
			in:     []byte(`{"name": "foo", "version>=": "1.3.37"}`),
			expDep: vcpkgDependency{Name: "foo", MinVersion: "1.3.37"},
		},
		{
			name:   "bare-string dependency",
			in:     []byte(`"bar"`),
			expDep: vcpkgDependency{Name: "bar"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var dep vcpkgDependency
			err := json.Unmarshal(tc.in, &dep)
			if err != nil {
				t.Fatal(err)
			}

			if tc.expDep.MinVersion != dep.MinVersion {
				t.Fatalf("expected MinVersion to be %s but got %s", tc.expDep.MinVersion, dep.MinVersion)
			}

			if tc.expDep.Name != dep.Name {
				t.Fatalf("expected Name to be %s but got %s", tc.expDep.Name, dep.Name)
			}
		})
	}
}
