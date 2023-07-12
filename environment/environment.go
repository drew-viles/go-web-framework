/*
Copyright 2023 Drew Viles.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package environment manages the environment variables passed into the applications
package environment

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
)

var apiSecret []byte

// GetAPISecret exists as we don't want others accessing the var for editing purposes.
func GetAPISecret() []byte {
	return apiSecret
}

// ReadEnvironmentFile reads the content of the web.config yaml file and parses them into a ConfigMap struct.
func ReadEnvironmentFile() (*ConfigMap, error) {
	log.Println("reading environment file")

	parseConfig()
	configMap := ConfigMap{
		App: web{
			FQDN:        viper.GetString("app.fqdn"),
			Env:         viper.GetString("app.env"),
			IP:          viper.GetString("app.ip"),
			Port:        viper.GetInt("app.port"),
			DomainShort: viper.GetString("app.domain_short"),
			SSL: certs{
				PrivateKey: viper.GetString("app.ssl.private_key"),
				PublicKey:  viper.GetString("app.ssl.public_key"),
				CAKey:      viper.GetString("app.ssl.ca_key"),
			},
		},
		Api: api{
			ApiSecret:   viper.GetString("api.api_secret"),
			ApiEndpoint: viper.GetString("api.api_endpoint"),
			Security: security{
				Token: token{
					Value:           viper.GetString("api.security.token.value"),
					ExpiryDate:      viper.GetDuration("api.security.token.expiry_date"),
					RefreshInterval: viper.GetDuration("api.security.token.refresh_interval"),
				},
			},
		},
		DB: db{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetInt("db.port"),
			Username: viper.GetString("db.username"),
			Password: viper.GetString("db.password"),
		},
		Stripe: configStripe{
			SecretKey:     viper.GetString("stripe.secret_key"),
			PublicKey:     viper.GetString("stripe.public_key"),
			WebhookSecret: viper.GetString("stripe.webhook_secret"),
			AccountID:     viper.GetString("stripe.account_id"),
		},
	}

	return &configMap, nil
}

// parseConfig will use viper to parse a config.yaml file
func parseConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/dcp-web/")
	viper.AddConfigPath("$HOME/dcp-web")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalln("the config file was not found in any of the valid locations - /etc/dcp-web/config.yaml, $HOME/dcp-web/config.yaml or ./config.yaml")
		} else {
			log.Fatalln("something went wrong reading the config file - please ensure it is valid YAML")
		}
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})

	viper.WatchConfig()
}
