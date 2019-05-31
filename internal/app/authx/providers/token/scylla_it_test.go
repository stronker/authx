package token

import (
	"github.com/nalej/authx/internal/app/authx/utils"
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

var _ = ginkgo.Describe("ScyllaTokenProvider", func() {

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
	sp := NewScyllaTokenProvider(scyllaHost, scyllaPort, nalejKeySpace)

	// disconnect
	ginkgo.AfterSuite(func() {
		sp.Disconnect()
	})

	TokenContexts(sp)

})
