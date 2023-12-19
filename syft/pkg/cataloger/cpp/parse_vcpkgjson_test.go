package cpp

import (
	"encoding/json"
	"testing"
)

func TestVcpkgDependencyUnmarshal(t *testing.T) {
	tests := []struct {
		name          string
		in            json.RawMessage
		expName       string
		expMinVersion string
	}{
		{
			name:          "json encoded dependency",
			in:            []byte(`{"name": "foo"}`),
			expName:       "foo",
			expMinVersion: "",
		},
		{
			name:          "json encoded dependency with MinVersion",
			in:            []byte(`{"name": "foo", "version>=": "1.3.37"}`),
			expName:       "foo",
			expMinVersion: "1.3.37",
		},
		{
			name:          "bare-string dependency",
			in:            []byte(`"bar"`),
			expName:       "bar",
			expMinVersion: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var dep vcpkgDependency
			err := json.Unmarshal(tc.in, &dep)
			if err != nil {
				t.Fatal(err)
			}

			if tc.expMinVersion != dep.MinVersion {
				t.Fatalf("expected MinVersion to be %s but got %s", tc.expMinVersion, dep.MinVersion)
			}

			if tc.expName != dep.Name {
				t.Fatalf("expected Name to be %s but got %s", tc.expName, dep.Name)
			}
		})
	}
}
