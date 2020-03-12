package config

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"time"
)

const (
	JWTDuration       = time.Hour * 24 * 365
	DefaultConfigPath = "./default_config.json"
	VerySecretKey     = "NKvqtjzgOP6adZMjmfc9tdhC7yB3UVcqMT89coA0lwdn8QpKHQ1uHCrSvX-w1jzD-PH1CFaKZD-iZiCohBBWlgpV_MnlfDdm-Wfn4aOb5dGgA1opsIVoQTa14T0eYPoEDQKsIc6tixba6gs_MPDIS0wSfKj4tLiYzYTooBcS-KR8YOKsykAGRz0dT-RDBJKXB0mEaMQ-GPpjw1wS2uj2_gAvIuhwvHtDRG1fUQ3RXgkAejF2oEDjhixxkKmbi7OZWpVt3xbozpBFpBpkJD7PqWhNmfTiFQb-aycP9NMcvneSZkIkH65_smsjZI4Ec10OrMOlGXlFR3jM88Ik8PRkeQ"
)

type PathConfig string

type Config struct {
	PostgresHost                  string `json:"postgresHost,omitempty"`
	PostgresPort                  uint32 `json:"postgresPort,omitempty"`
	PostgresPassword              string `json:"postgresPassword,omitempty"`
	PostgresSSLMode               string `json:"postgresSSLMode,omitempty"`
	PostgresDBName                string `json:"postgresDBName,omitempty"`
	PostgresUserName              string `json:"postgresUserName,omitempty"`
	JWTIdentityKey                string `json:"jWTIdentityKey,omitempty"`
	MediaDir                      string `json:"mediaDir,omitempty"`
	DefaultGroupImagePath         string `json:"defaultGroupImagePath,omitempty"`
	DefaultOrgImagePath           string `json:"defaultOrgImagePath,omitempty"`
	DefaultUserImagePath          string `json:"defaultUserImagePath,omitempty"`
	Version                       string `json:"version,omitempty"`
	ApiFnsKey                     string `json:"apiFnsKey,omitempty"`
	Host                          string `json:"host,omitempty"`
	Port                          uint32 `json:"port,omitempty"`
	EventEmitterChannelBufferSize uint32 `json:"eventEmitterChannelBufferSize,omitempty"`
	ConnectionManagerBufferSize   uint32 `json:"connectionManagerBufferSize,omitempty"`
	RatingEventDuration           uint32 `json:"ratingEventDuration,omitempty"`
	ElasticSearchUrl              string `json:"elasticSearchUrl,omitempty"`
	ZomboDBOn                     bool   `json:"zomboDBOn,omitempty"`
	DebugMode                     bool   `json:"debugMode,omitempty"`
	RatingEventDebugDuration      uint32 `json:"ratingEventDebugDuration,omitempty"`
	SigningAlgorithm              string `json:"signingAlgorithm"`
	BugReportsDir                 string `json:"bugReportsDir"`
	RatingFile                    string `json:"ratingFile"`
}

type notDefaultConfig struct {
	PostgresHost                  *string `json:"postgresHost,omitempty"`
	PostgresPort                  *uint32 `json:"postgresPort,omitempty"`
	PostgresPassword              *string `json:"postgresPassword,omitempty"`
	PostgresSSLMode               *string `json:"postgresSSLMode,omitempty"`
	PostgresDBName                *string `json:"postgresDBName,omitempty"`
	PostgresUserName              *string `json:"postgresUserName,omitempty"`
	JWTIdentityKey                *string `json:"jWTIdentityKey,omitempty"`
	MediaDir                      *string `json:"mediaDir,omitempty"`
	DefaultGroupImagePath         *string `json:"defaultGroupImagePath,omitempty"`
	DefaultOrgImagePath           *string `json:"defaultOrgImagePath,omitempty"`
	DefaultUserImagePath          *string `json:"defaultUserImagePath,omitempty"`
	Version                       *string `json:"version,omitempty"`
	ApiFnsKey                     *string `json:"apiFnsKey,omitempty"`
	Host                          *string `json:"host,omitempty"`
	Port                          *uint32 `json:"port,omitempty"`
	EventEmitterChannelBufferSize *uint32 `json:"eventEmitterChannelBufferSize,omitempty"`
	ConnectionManagerBufferSize   *uint32 `json:"connectionManagerBufferSize,omitempty"`
	RatingEventDuration           *uint32 `json:"ratingEventDuration"`
	ElasticSearchUrl              *string `json:"elasticSearchUrl"`
	ZomboDBOn                     *bool   `json:"zomboDBOn,omitempty"`
	DebugMode                     *bool   `json:"debugMode,omitempty"`
	RatingEventDebugDuration      *uint32 `json:"ratingEventDebugDuration,omitempty"`
	SigningAlgorithm              *string `json:"signingAlgorithm"`
	BugReportsDir                 *string `json:"bugReportsDir"`
	RatingFile                    *string `json:"ratingFile"`
}

func copyIfNotNil(src *notDefaultConfig, dst *Config) {
	v1 := reflect.ValueOf(dst).Elem()
	v2 := reflect.ValueOf(src).Elem()
	t1 := reflect.TypeOf(dst).Elem()
	t2 := reflect.TypeOf(src).Elem()
	for i := 0; i < t1.NumField(); i++ {
		if isSame(t1.Field(i), t2.Field(i)) && !v2.Field(i).IsNil() {
			v1.Field(i).Set(v2.Field(i).Elem())
		}
	}
}

func isSame(notPtr, ptr reflect.StructField) bool {
	return notPtr.Type == ptr.Type.Elem() && notPtr.Name == ptr.Name
}

func ParseConfig(path PathConfig) (*Config, error) {
	defaultConfigBytes, err := ioutil.ReadFile(DefaultConfigPath)
	if err != nil {
		return nil, err
	}

	var defaultConfig Config
	err = json.Unmarshal(defaultConfigBytes, &defaultConfig)
	if err != nil {
		return nil, err
	}

	if len(path) > 0 {
		configBytes, err := ioutil.ReadFile(string(path))
		if err != nil {
			return nil, err
		}

		var config notDefaultConfig
		err = json.Unmarshal(configBytes, &config)
		if err != nil {
			return nil, err
		}

		copyIfNotNil(&config, &defaultConfig)
	}

	return &defaultConfig, nil
}
