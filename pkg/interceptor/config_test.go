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

package interceptor

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

const ValidConfig = `
{
	"allows_all":true,
	"permissions": {
		"/authx.Authx/AddBasicCredentials":{
			"must": ["primitive1"]
		}
	}
}
`

const InValidConfig = `
{
	"allows_all":true,
	"permissions": {
		"must": ["primitive1"]
	}
}
`

var _ = ginkgo.Describe("Load config", func() {
	ginkgo.Context("with a valid config file", func() {

		ginkgo.It("should add basic credentials with correct roleID and correct JWT", func() {
			var tmpFile *os.File
			tmpFile, fileErr := ioutil.TempFile("", "load-test")
			gomega.Expect(fileErr).To(gomega.Succeed())
			gomega.Expect(tmpFile).NotTo(gomega.BeNil())
			tmpFile.WriteString(ValidConfig)

			cfg, err := LoadAuthorizationConfig(tmpFile.Name())
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(cfg).NotTo(gomega.BeNil())
			gomega.Expect(cfg.AllowsAll).To(gomega.BeTrue())
			gomega.Expect(cfg.Permissions).To(gomega.HaveKey("/authx.Authx/AddBasicCredentials"))
		})

	})

	ginkgo.Context("with a invalid path", func() {

		ginkgo.It("should fail", func() {
			cfg, err := LoadAuthorizationConfig("/invalidPath/test.txt")
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(cfg).To(gomega.BeNil())
		})
	})

	ginkgo.Context("with a bad config file", func() {

		ginkgo.It("should fail", func() {
			var tmpFile *os.File
			tmpFile, fileErr := ioutil.TempFile("", "load-test")
			gomega.Expect(fileErr).To(gomega.Succeed())
			gomega.Expect(tmpFile).NotTo(gomega.BeNil())
			tmpFile.WriteString(InValidConfig)

			cfg, err := LoadAuthorizationConfig(tmpFile.Name())
			gomega.Expect(err).To(gomega.HaveOccurred())
			gomega.Expect(cfg).To(gomega.BeNil())
		})
	})
})
