/*
 * Ceph Nano (C) 2018 Red Hat, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

/*
 * Below main package has canonical imports for 'go get' and 'go build'
 * to work with all other clones of github.com/ceph/cn repository. For
 * more information refer https://golang.org/doc/go1.4#canonicalimports
 */

package cmd

import (
	"path"

	"github.com/alecthomas/units"
	"github.com/spf13/viper"
)

func readConfigFile(customFile ...string) {
	setDefaultConfig()
	if len(customFile) > 0 {
		var filename = path.Base(customFile[0])
		var fileDir = path.Dir(customFile[0])
		viper.SetConfigFile(filename)
		viper.AddConfigPath(fileDir)
		err := viper.ReadInConfig() // Find and read the config file
		// If there is no configuration file, that's an error
		if err != nil {
			panic(err)
		}
	} else {
		viper.SetConfigName("cn")         // name of config file (without extension)
		viper.AddConfigPath("/etc/ceph/") // path to look for the config file in
		viper.AddConfigPath("$HOME/.cn/") // call multiple times to add many search paths
		viper.AddConfigPath(".")          // optionally look for config in the working directory
		viper.ReadInConfig()              // Find and read the config file, we don't really care if no config file is found
	}
}

// Set the default values for defined types
// If the configuration file is missing, this section will generated the mandatory elements
func setDefaultConfig() {
	viper.SetDefault("default.use_default", "true") // All containers inherit from default
	viper.SetDefault("default.MemorySize", "512MB")
}

func getStringFromConfig(name string, containerName string) string {
	var value = ""
	// If we are requested to get the status of use_default, we cannot call useDefault ;)
	if name == "use_default" || useDefault(containerName) {
		value = viper.GetString("default" + "." + name)
	}
	containerValue := viper.GetString(containerName + "." + name)
	if len(containerValue) > 0 {
		value = containerValue
	}
	return value
}

func useDefault(containerName string) bool {
	useDefaultValue := getStringFromConfig("use_default", containerName)
	if (len(useDefaultValue) > 0) && (useDefaultValue == "true") {
		return true
	}
	return false
}

func getStringMapFromConfig(name string, containerName string) map[string]interface{} {
	var defaultConfig = make(map[string]interface{})
	if useDefault(containerName) {
		defaultConfig = viper.GetStringMap("default" + "." + name)
	}
	containerValues := viper.GetStringMap(containerName + "." + name)
	if len(containerValues) > 0 {
		for key, value := range containerValues {
			defaultConfig[key] = value
		}
	}
	return defaultConfig
}

func getMemorySize(containerName string) int64 {
	var bytes units.Base2Bytes
	var err error
	bytes, err = units.ParseBase2Bytes(getStringFromConfig("MemorySize", containerName))
	if err != nil {
		panic(err)
	}
	return int64(bytes)
}

func getCephConf(containerName string) map[string]interface{} {
	return getStringMapFromConfig("ceph.conf", containerName)
}