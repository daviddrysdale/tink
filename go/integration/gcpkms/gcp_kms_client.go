// Copyright 2017 Google Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
////////////////////////////////////////////////////////////////////////////////

package gcpkms

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"google3/third_party/golang/google_api/cloudkms/v1/cloudkms"
	"google3/third_party/golang/oauth2/google/google"
	"google3/third_party/golang/oauth2/oauth2"
	"github.com/google/tink/go/tink"
)

const (
	gcpPrefix = "gcp-kms://"
)

var (
	errCred = errors.New("invalid credential path")
)

// GCPClient represents a client that connects to the GCP KMS backend.
type GCPClient struct {
	keyURI string
	kms    *cloudkms.Service
}

var _ tink.KMSClient = (*GCPClient)(nil)

// NewGCPClient returns a new client to GCP KMS. It does not have an established session.
func NewGCPClient(URI string) (*GCPClient, error) {
	if !strings.HasPrefix(strings.ToLower(URI), gcpPrefix) {
		return nil, errors.New(fmt.Sprintf("key URI must start with %s", gcpPrefix))
	}

	return &GCPClient{
		keyURI: URI,
	}, nil

}

// Supported true if this client does support keyURI
func (g *GCPClient) Supported(keyURI string) bool {
	if (len(g.keyURI) > 0) && (strings.Compare(strings.ToLower(g.keyURI), strings.ToLower(keyURI)) == 0) {
		return true
	}
	return ((len(g.keyURI) == 0) && (strings.HasPrefix(strings.ToLower(keyURI), gcpPrefix)))
}

// LoadCredentials loads the credentials in credentialPath. If credentialPath is  null, loads the
// default credentials.
func (g *GCPClient) LoadCredentials(credentialPath string) (interface{}, error) {
	ctx := context.Background()
	if len(credentialPath) <= 0 {
		return nil, errCred
	}
	data, err := ioutil.ReadFile(credentialPath)
	if err != nil {
		return nil, errCred
	}
	creds, err := google.CredentialsFromJSON(ctx, data, "https://www.googleapis.com/auth/cloudkms")
	if err != nil {
		return nil, errCred
	}
	client := oauth2.NewClient(ctx, creds.TokenSource)
	kmsService, err := cloudkms.New(client)
	if err != nil {
		return nil, err
	}
	g.kms = kmsService
	return g, nil
}

// LoadDefaultCredentials loads with the default credentials.
func (g *GCPClient) LoadDefaultCredentials() (interface{}, error) {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		return nil, err
	}

	kmsService, err := cloudkms.New(client)
	if err != nil {
		return nil, err
	}
	g.kms = kmsService
	return g, nil
}

// GetAead  gets an Aead backend by keyURI.
func (g *GCPClient) GetAead(keyURI string) (tink.AEAD, error) {
	if len(g.keyURI) > 0 && strings.Compare(strings.ToLower(g.keyURI), strings.ToLower(keyURI)) != 0 {
		return nil, errors.New(fmt.Sprintf("this client is bound to %s, cannot load keys bound to %s", g.keyURI, keyURI))
	}
	uri, err := validateTrimKMSPrefix(g.keyURI, gcpPrefix)
	if err != nil {
		return nil, err
	}
	return NewGCPAead(uri, g.kms), nil
}

func validateKMSPrefix(keyURI, prefix string) bool {
	if len(keyURI) > 0 && strings.HasPrefix(strings.ToLower(keyURI), gcpPrefix) {
		return true
	}
	return false
}

func validateTrimKMSPrefix(keyURI, prefix string) (string, error) {
	if !validateKMSPrefix(keyURI, prefix) {
		return "", errors.New(fmt.Sprintf("key URI must start with %s", prefix))
	}
	return strings.TrimPrefix(keyURI, prefix), nil
}