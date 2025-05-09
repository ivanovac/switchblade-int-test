package integration_test

import (
	"testing"

	"github.com/cloudfoundry/switchblade"

	. "github.com/cloudfoundry/switchblade/matchers"
	. "github.com/onsi/gomega"
)

func TestDocker(t *testing.T) {
	var (
		Expect     = NewWithT(t).Expect
		Eventually = NewWithT(t).Eventually
	)

	// Create an instance of a Docker platform. A GitHub token is required to
	// make API requests to GitHub fetching buildpack details.
	platform, err := switchblade.NewPlatform(switchblade.Docker, "use-personal-access-tokens-for-github", "cflinuxfs4") //
	Expect(err).NotTo(HaveOccurred())

	// Deploy an application called "my-app" onto Docker with source code
	// located at /path/to/my/app/source. This is similar to the following `cf`
	// command, but running locally on your Docker daemon:
	//   cf push my-app -p /path/to/my/app
	deployment, logs, err := platform.Deploy.Execute("test-app", "path-to-test-app")
	Expect(err).NotTo(HaveOccurred())

	// Assert that the deployment logs contain a line that contains the substring
	// "Installing dependency..."
	Expect(logs).To(ContainLines(ContainSubstring("Installing dependency...")))

	// Assert that the deployment results in an application instance that serves
	// "Hello, world!" over HTTP.
	Eventually(deployment).Should(Serve(ContainSubstring("Hello, world!")))

	// Delete the application from the platform.
	Expect(platform.Delete.Execute("test-app")).To(Succeed())
}
