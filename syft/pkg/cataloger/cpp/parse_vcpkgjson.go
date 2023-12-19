package cpp

import (
	"encoding/json"
	"io"
	"strings"

	"github.com/anchore/syft/syft/artifact"
	"github.com/anchore/syft/syft/file"
	"github.com/anchore/syft/syft/pkg"
	"github.com/anchore/syft/syft/pkg/cataloger/generic"
)

type vcpkgJSON struct {
	DefaultFeatures []defaultFeatureSetting    `json:"default-features"`
	Features        map[string]feature         `json:"features"`
	Dependencies    map[string]vcpkgDependency `json:"dependencies"`
	Overrides       map[string]string          `json:"overrides"`
}

// Confusinglu, vcpkg has twp different concepts, both called Features in their documentation:
// 1) Features which contains data, and potentially dependencies, for features that are optionally enabled by the user when building the application.
//   - In vcppkg docs, this is same datastructure called a Feature Object
//
// 2) Features which only tell you if features (1) are enabled by default or not. In vcpkg docs this is called a Feature
//
// To clarify the distinction in this code, features (1) are called features, and features (2) are called defaultFeatureSettings
type feature struct {
	Description  string            `json:"description"` // Description is the only required field in a feature
	Dependencies []vcpkgDependency `json:"dependencies"`
	Supports     []string          `json:"supports"`
	License      *string           `json:"license"` // License is possibly null, so to use unmarshal, we need to use a pointer
}

type defaultFeatureSetting struct {
	Name     string `json:"name"`
	Platform string `json:"platform"`
}

func (f *defaultFeatureSetting) UnmarshalJSON(b []byte) error {
	// The structure of a feature is either a string with the name of the feature, or it's a json object with the name and the platform
	var s string
	if err := json.Unmarshal(b, &s); err == nil {
		f.Name = s
		return nil
	}

	return json.Unmarshal(b, f)
}

type vcpkgDependency struct {
	Name            string                  `json:"name"`
	MinVersion      string                  `json:"version>="`
	DefaultFeatures bool                    `json:"default-features"`
	Features        []defaultFeatureSetting `json:"features"`
	Platform        string                  `json:"platform"`
}

func (v *vcpkgDependency) UnmarshalJSON(b []byte) error {
	// The structure of a dependency is either a string with the name of the dependency, or it's a json object with the name, and possibly the version

	var dep struct {
		Name       string `json:"name"`
		MinVersion string `json:"version>="`
	}

	if err := json.Unmarshal(b, &dep); err == nil {
		v.Name = dep.Name
		v.MinVersion = dep.MinVersion
		return nil
	}

	// The raw string representation of the dependency is a quoted string, so we need to remove the quotes
	v.Name = strings.ReplaceAll(string(b), `"`, ``)
	return nil
}

func parseVcpkgJSON(_ file.Resolver, _ *generic.Environment, reader file.LocationReadCloser) ([]pkg.Package, []artifact.Relationship, error) {
	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, nil, err
	}

	var v vcpkgJSON
	if err := json.Unmarshal(bytes, &v); err != nil {
		return nil, nil, err
	}

	return nil, nil, nil
}
