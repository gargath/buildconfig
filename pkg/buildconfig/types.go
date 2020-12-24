package buildconfig

type options struct {
	info BuildConfig
}

// BuildConfig holds parsed config data from buildconfig file
type BuildConfig struct {
	Binary       string
	Dependencies []struct {
		Name         string       `yaml:"name"`
		VersionCheck VersionCheck `yaml:"version_check,omitempty"`
	} `yaml:"dependencies"`
}

// VersionCheck holds version check data from buildconfig file
type VersionCheck struct {
	Command    string `yaml:"command"`
	MinVersion string `yaml:"min_version,omitempty"`
}
