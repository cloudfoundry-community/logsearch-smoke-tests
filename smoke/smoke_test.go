package smoke

import (
	"fmt"
	"strings"
	"time"

	"github.com/pborman/uuid"

	"github.com/pivotal-cf-experimental/cf-test-helpers/cf"
	"github.com/pivotal-cf-experimental/cf-test-helpers/runner"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"
)

var _ = Describe("Logsearch", func() {
	var timeout = time.Second * 60
	var appPath = "../assets/logsearch-example-app"

	var appName string
	var elasticEndpoint string

	randomName := func() string {
		return uuid.NewRandom().String()
	}

	elaticUri := func(elasticEndpoint string) string {
		return "http://" + config.ElasticsearchMasterIpAddress + ":" + config.ElasticsearchMasterPort
	}

	BeforeEach(func() {
		appName = randomName()
		Eventually(cf.Cf("push", appName, "-m", "128M", "-p", appPath, "-u", "none", "-no-start"), config.ScaledTimeout(timeout)).Should(Exit(0))
	})

	AfterEach(func() {
		Eventually(cf.Cf("delete", appName, "-f"), config.ScaledTimeout(timeout)).Should(Exit(0))
	})

	It("can see app messages in the elasticsearch", func() {
		Eventually(cf.Cf("start", appName), 5*60*time.Second).Should(Exit(0))

		current_time := time.Now().Local()

		testUri := fmt.Sprintf(
			elaticUri(elasticEndpoint)+"/"+config.ElasticsearchAppIndex+
				"-"+ctx.RegularUserContext().Org+
				"-"+ctx.RegularUserContext().Space+
				"-"+"%v"+"/_search?q=@cf.app:"+"%s",
			current_time.Format("2006.01.02"), appName)

		fmt.Println("Curling url: ", testUri)

		curl := runner.Curl(strings.ToLower(testUri)).Wait(timeout)
		Expect(curl).To(Exit(0))
		elasticResponse := string(curl.Out.Contents())

		Eventually(elasticResponse).Should(ContainSubstring(appName))
		fmt.Println("\n")
	})

})
