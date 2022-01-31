/*
Copyright 2017 Frederic Branczyk All rights reserved.

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

package authn

import (
	"io/ioutil"

	"k8s.io/apiserver/pkg/authentication/authenticator"
	"k8s.io/apiserver/pkg/authentication/request/bearertoken"
	"k8s.io/apiserver/plugin/pkg/authenticator/token/oidc"
)

// OIDCConfig represents configuration used for JWT request authentication
type OIDCConfig struct {
	IssuerURL            string
	ClientID             string
	CAFile               string
	UsernameClaim        string
	UsernamePrefix       string
	GroupsClaim          string
	GroupsPrefix         string
	SupportedSigningAlgs []string
}

type caContentProvider func() []byte

func (f caContentProvider) CurrentCABundleContent() []byte {
	return f()
}

func newCAContentProvider(filename string) (caContentProvider, error) {
	pemBlock, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return func() []byte {
		return pemBlock
	}, nil
}

// NewOIDCAuthenticator returns OIDC authenticator
func NewOIDCAuthenticator(config *OIDCConfig) (authenticator.Request, error) {
	cacp, err := newCAContentProvider(config.CAFile)
	if err != nil {
		return nil, err
	}

	tokenAuthenticator, err := oidc.New(oidc.Options{
		IssuerURL:            config.IssuerURL,
		ClientID:             config.ClientID,
		CAContentProvider:    cacp,
		UsernameClaim:        config.UsernameClaim,
		UsernamePrefix:       config.UsernamePrefix,
		GroupsClaim:          config.GroupsClaim,
		GroupsPrefix:         config.GroupsPrefix,
		SupportedSigningAlgs: config.SupportedSigningAlgs,
	})
	if err != nil {
		return nil, err
	}

	return bearertoken.New(tokenAuthenticator), nil
}
