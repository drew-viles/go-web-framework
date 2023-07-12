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

package environment

import (
	"time"
)

type web struct {
	FQDN        string `yaml:"fqdn" validate:"omitempty"`
	Env         string `yaml:"env" validate:"omitempty"`
	IP          string `yaml:"ip" validate:"omitempty"`
	Port        int    `yaml:"port" validate:"omitempty"`
	DomainShort string `yaml:"domain_short" validate:"omitempty"`
	SSL         certs
}

type api struct {
	ApiSecret   string `yaml:"api_secret" validate:"omitempty"`
	ApiEndpoint string `yaml:"api_endpoint" validate:"omitempty"`
	Security    security
}

type security struct {
	Token token
}

type certs struct {
	PrivateKey string `yaml:"private_key" validate:"omitempty"`
	PublicKey  string `yaml:"public_key" validate:"omitempty"`
	CAKey      string `yaml:"ca_key" validate:"omitempty"`
}

type token struct {
	Value           string        `yaml:"value" validate:"omitempty"`
	ExpiryDate      time.Duration `yaml:"expiry_time" validate:"omitempty"`
	RefreshInterval time.Duration `yaml:"refresh_interval" validate:"omitempty"`
}

type db struct {
	Host     string `yaml:"host" validate:"omitempty"`
	Port     int    `yaml:"port" validate:"omitempty"`
	Username string `yaml:"username" validate:"omitempty"`
	Password string `yaml:"password" validate:"omitempty"`
}

type configStripe struct {
	SecretKey     string `yaml:"secret_key" validate:"omitempty"`
	PublicKey     string `yaml:"public_key" validate:"omitempty"`
	WebhookSecret string `yaml:"webhook_secret" validate:"omitempty"`
	AccountID     string `yaml:"account_id" validate:"omitempty"`
}

type ConfigMap struct {
	App    web          `yaml:"web"`
	Api    api          `yaml:"api"`
	DB     db           `yaml:"db"`
	Stripe configStripe `yaml:"stripe"`
}
