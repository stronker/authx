/*
 * Copyright 2019 Nalej
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
 *
 */

package certificates

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"github.com/nalej/grpc-authx-go"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"github.com/stronker/authx/internal/app/authx/config"
	"math/big"
	"time"
)

func createTestCA() (*x509.Certificate, *rsa.PrivateKey) {
	
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	gomega.Expect(err).To(gomega.Succeed())
	
	caCert := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Issuer: pkix.Name{
			Organization: []string{"Nalej"},
		},
		Subject: pkix.Name{
			Organization: []string{"Nalej"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(CertValidity),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            0,
		MaxPathLenZero:        true,
		DNSNames:              []string{"*.fake.nalej.tech"},
	}
	publicKey := &privateKey.PublicKey
	rawCert, err := x509.CreateCertificate(rand.Reader, &caCert, &caCert, publicKey, privateKey)
	gomega.Expect(err).To(gomega.Succeed())
	
	cert, err := x509.ParseCertificate(rawCert)
	gomega.Expect(err).To(gomega.Succeed())
	
	return cert, privateKey
}

var testManager Manager

var _ = ginkgo.Describe("With a manager", func() {
	ginkgo.BeforeSuite(func() {
		testCA, testPK := createTestCA()
		helper := &CertHelper{
			CACert:     testCA,
			PrivateKey: testPK,
		}
		
		emptyCfg := config.Config{}
		
		testManager = NewManager(emptyCfg, helper)
	})
	
	ginkgo.It("should be able to generate an edge controller certificate", func() {
		request := &grpc_authx_go.EdgeControllerCertRequest{
			OrganizationId:   "organization_id",
			EdgeControllerId: "edge_controller_id",
			Name:             "Fake EC",
			Ips:              []string{"10.0.0.10", "192.168.250.10"},
		}
		ecCert, err := testManager.CreateControllerCert(request)
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(ecCert).ToNot(gomega.BeNil())
		gomega.Expect(ecCert.Certificate).ShouldNot(gomega.BeNil())
		gomega.Expect(ecCert.PrivateKey).ShouldNot(gomega.BeNil())
		
		x509Cert, cErr := tls.X509KeyPair([]byte(ecCert.Certificate), []byte(ecCert.PrivateKey))
		gomega.Expect(cErr).To(gomega.Succeed())
		
		block, _ := pem.Decode([]byte(ecCert.Certificate))
		gomega.Expect(block).ShouldNot(gomega.BeNil())
		
		cert, perr := x509.ParseCertificate(block.Bytes)
		gomega.Expect(perr).To(gomega.Succeed())
		
		log.Info().Interface("cert", x509Cert).Msg("result")
		log.Info().Interface("cert", cert).Msg("result")
	})
	
})
