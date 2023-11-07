package falcon

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	"log"
	"os"
	"regexp"
	"strings"
)

type FileType string

const (
	Env             FileType = ".env"
	ApplicationYaml          = "application.yaml"
)

type Falcon struct {
	resourcePath string
}

type falconBuilder struct {
	falcon Falcon
}

func NewFalconBuilder() *falconBuilder {
	return &falconBuilder{}
}

func (f *falconBuilder) ResourcePath(resourcePath string) *falconBuilder {
	f.falcon.resourcePath = resourcePath
	return f
}

func (f *falconBuilder) Build() Service {
	fmt.Println("rootPath -> ", f.falcon.resourcePath)

	if isEmpty(f.falcon.resourcePath) {
		panic("resource path is required to load properties")
	}
	return &Falcon{
		resourcePath: f.falcon.resourcePath,
	}
}

func (f *Falcon) GetByKey(key string) string {
	property := os.Getenv(key)
	if isEmpty(property) {
		property = viper.GetString(key)
		if containsSoEnv(property) {
			match := getValueByExpression(property)
			if strings.Contains(match, ":") {
				splitEnv := strings.SplitN(match, ":", 2)
				envValue := os.Getenv(splitEnv[0])
				if isEmpty(envValue) {
					return splitEnv[1]
				}
				return envValue
			} else {
				envValue := os.Getenv(match)
				if isEmpty(envValue) {
					return match
				}
				return envValue
			}
		}
	}
	return property
}

func getValueByExpression(property string) string {
	regex := regexp.MustCompile("\\${(.*?)}")
	match := regex.FindStringSubmatch(property)
	return match[1]
}

func containsSoEnv(property string) bool {
	if len(property) < 1 {
		return false
	}
	firstProp := property[0:2]
	lastProp := property[len(property)-1:]
	return firstProp == "${" && lastProp == "}"
}

func (f *Falcon) loadEnv() {
	err := gotenv.Load(fmt.Sprintf("%s/%s", f.resourcePath, ".env"))
	if err != nil {
		log.Printf("Error loading %s file", fmt.Sprintf("%s/%s", f.resourcePath, ".env"))
	}
}

func (f *Falcon) configViper() {
	viper.AddConfigPath(f.resourcePath)
	viper.SetConfigName("application")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix("")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s", err)
	}
}

func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
