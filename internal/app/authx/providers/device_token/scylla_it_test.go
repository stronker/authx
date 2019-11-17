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

package device_token

import (
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
	"github.com/stronker/authx/internal/app/authx/utils"
	"os"
	"strconv"
)

var _ = ginkgo.Describe("ScyllaDeviceTokenProvider", func() {
	
	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}
	
	var scyllaHost = os.Getenv("IT_SCYLLA_HOST")
	if scyllaHost == "" {
		ginkgo.Fail("missing environment variables")
	}
	
	scyllaPort, _ := strconv.Atoi(os.Getenv("IT_SCYLLA_PORT"))
	
	if scyllaPort <= 0 {
		ginkgo.Fail("missing environment variables")
	}
	
	var nalejKeySpace = os.Getenv("IT_NALEJ_KEYSPACE")
	if nalejKeySpace == "" {
		ginkgo.Fail("missing environment variables")
		
	}
	
	// create a provider and connect it
	sp := NewScyllaDeviceTokenProvider(scyllaHost, scyllaPort, nalejKeySpace)
	
	// disconnect
	ginkgo.AfterSuite(func() {
		sp.Disconnect()
	})
	
	DeviceTokenContexts(sp)
	
})
