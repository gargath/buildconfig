package buildconfig

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

func parseInfo(buildconfigPath string) (*options, error) {
	buildconfigFile, err := ioutil.ReadFile(buildconfigPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read buildconfig file %s: %v", buildconfigPath, err)
	}
	var buildconfigYaml BuildConfig
	err = yaml.Unmarshal(buildconfigFile, &buildconfigYaml)
	if err != nil {
		return nil, fmt.Errorf(" buildconfig file %s is invald: %v", buildconfigPath, err)
	}
	opts := &options{
		info: buildconfigYaml,
	}
	return opts, nil
}
