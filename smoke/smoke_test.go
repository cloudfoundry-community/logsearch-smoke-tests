package smoke

import (
	"bytes"
	"fmt"
	"html/template"
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

	elasticURI := func(elasticEndpoint string) string {
		return "http://" + config.ElasticsearchMasterIpAddress + ":" + config.ElasticsearchMasterPort
	}

	elasticIndex := func(raw string) (string, error) {
		tpl, err := template.New("index").Parse(raw)
		var buf bytes.Buffer
		if err != nil {
			return "", err
		}
		if err := tpl.Execute(&buf, map[string]interface{}{
			"Org":   ctx.RegularUserContext().Org,
			"Space": ctx.RegularUserContext().Space,
			"Time":  time.Now().Local(),
		}); err != nil {
			return "", err
		}
		return buf.String(), nil
	}

	BeforeEach(func() {
		appName = randomName()
		Eventually(cf.Cf("push", appName, "-m", "128M", "-p", appPath, "-u", "none", "--no-start"), config.ScaledTimeout(timeout)).Should(Exit(0))
	})

	AfterEach(func() {
		Eventually(cf.Cf("delete", appName, "-f"), config.ScaledTimeout(timeout)).Should(Exit(0))
	})

	It("can see app messages in the elasticsearch", func() {
		Eventually(cf.Cf("start", appName), 5*60*time.Second).Should(Exit(0))

		index, err := elasticIndex(config.ElasticsearchAppIndex)
		Expect(err).NotTo(HaveOccurred())

		testURI := elasticURI(elasticEndpoint) + "/" + index + "/_search?q=@cf.app:" + appName
		fmt.Println("Curling url: ", testURI)

		curl := runner.Curl(strings.ToLower(testURI)).Wait(timeout)
		Expect(curl).To(Exit(0))
		elasticResponse := string(curl.Out.Contents())

		Eventually(elasticResponse).Should(ContainSubstring(appName))
		fmt.Println()
	})
})
