package device

/*
create KEYSPACE authx WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1};

create table IF NOT EXISTS authx.deviceCredentials (organization_id text, device_group_id text, device_id text, device_api_key text, enabled boolean, PRIMARY KEY ((organization_id, device_group_id), device_id));
create table IF NOT EXISTS authx.deviceGroupCredentials (organization_id text, device_group_id text,  device_group_api_key text, enabled boolean, default_device_connectivity boolean, PRIMARY KEY (organization_id, device_group_id));

create INDEX IF NOT EXISTS device_group_api ON authx.devicegroupcredentials ( device_group_api_key);
create INDEX IF NOT EXISTS device_api ON authx.devicecredentials ( device_api_key);

IT_SCYLLA_HOST=127.0.0.1
RUN_INTEGRATION_TEST=true
IT_NALEJ_KEYSPACE=authx
IT_SCYLLA_PORT=9042

 */


import (
	"github.com/nalej/authx/internal/app/authx/utils"
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

var _ = ginkgo.Describe("Scylla cluster provider", func() {

	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}
	var scyllaHost = os.Getenv("IT_SCYLLA_HOST")
	if scyllaHost == "" {
		ginkgo.Fail("missing environment variables")
	}
	var nalejKeySpace = os.Getenv("IT_NALEJ_KEYSPACE")
	if scyllaHost == "" {
		ginkgo.Fail("missing environment variables")
	}
	scyllaPort, _ := strconv.Atoi(os.Getenv("IT_SCYLLA_PORT"))
	if scyllaPort <= 0 {
		ginkgo.Fail("missing environment variables")
	}

	// create a provider and connect it
	sp := NewScyllaDeviceCredentialsProvider(scyllaHost, scyllaPort, nalejKeySpace)

	ginkgo.AfterSuite(func() {
		sp.Disconnect()
	})

	RunTest(sp)

})
