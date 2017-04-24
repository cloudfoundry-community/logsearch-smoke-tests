package smoke

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/pivotal-cf-experimental/cf-test-helpers/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type testConfig struct {
	services.Config

	ElasticsearchMasterIpAddress string `json:"elasticsearch_master_ip"`
	ElasticsearchMasterPort      string `json:"elasticsearch_master_port"`
	ElasticsearchAppIndex        string `json:"elasticsearch_app_index"`
}

func loadConfig(path string) (cfg testConfig) {
	configFile, err := os.Open(path)
	if err != nil {
		fatal(err)
	}

	decoder := json.NewDecoder(configFile)
	if err = decoder.Decode(&cfg); err != nil {
		fatal(err)
	}

	return
}

var (
	config = loadConfig(os.Getenv("CONFIG"))
	ctx    services.Context
)

func fatal(err error) {
	fmt.Printf("ERROR: %s\n", err.Error())
	os.Exit(1)
}

func TestService(t *testing.T) {
	if err := services.ValidateConfig(&config.Config); err != nil {
		fatal(err)
	}

	ctx = services.NewContext(config.Config, "logsearch-smoke-tests")

	RegisterFailHandler(Fail)

	RunSpecs(t, "Logsearch-for-CloudFoundry Smoke Tests")
}

var _ = BeforeEach(func() {
	ctx.Setup()
})

var _ = AfterEach(func() {
	ctx.Teardown()
})
