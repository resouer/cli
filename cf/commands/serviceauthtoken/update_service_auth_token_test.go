package serviceauthtoken_test

import (
	"github.com/cloudfoundry/cli/cf/configuration"
	"github.com/cloudfoundry/cli/cf/models"
	testapi "github.com/cloudfoundry/cli/testhelpers/api"
	testcmd "github.com/cloudfoundry/cli/testhelpers/commands"
	testconfig "github.com/cloudfoundry/cli/testhelpers/configuration"
	testreq "github.com/cloudfoundry/cli/testhelpers/requirements"
	testterm "github.com/cloudfoundry/cli/testhelpers/terminal"

	. "github.com/cloudfoundry/cli/cf/commands/serviceauthtoken"
	. "github.com/cloudfoundry/cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("update-service-auth-token command", func() {
	var (
		ui                  *testterm.FakeUI
		configRepo          configuration.ReadWriter
		authTokenRepo       *testapi.FakeAuthTokenRepo
		requirementsFactory *testreq.FakeReqFactory
	)

	BeforeEach(func() {
		ui = &testterm.FakeUI{Inputs: []string{"y"}}
		authTokenRepo = &testapi.FakeAuthTokenRepo{}
		configRepo = testconfig.NewRepositoryWithDefaults()
		requirementsFactory = &testreq.FakeReqFactory{}
	})

	runCommand := func(args ...string) {
		testcmd.RunCommand(NewUpdateServiceAuthToken(ui, configRepo, authTokenRepo), args, requirementsFactory)
	}

	Describe("requirements", func() {
		It("fails with usage when not provided exactly three args", func() {
			requirementsFactory.LoginSuccess = true
			runCommand("some-token-label", "a-provider")
			Expect(ui.FailedWithUsage).To(BeTrue())
		})

		It("fails when not logged in", func() {
			runCommand("label", "provider", "token")
			Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
		})
	})

	Context("when logged in and the service auth token exists", func() {
		BeforeEach(func() {
			requirementsFactory.LoginSuccess = true
			foundAuthToken := models.ServiceAuthTokenFields{}
			foundAuthToken.Guid = "found-auth-token-guid"
			foundAuthToken.Label = "found label"
			foundAuthToken.Provider = "found provider"
			authTokenRepo.FindByLabelAndProviderServiceAuthTokenFields = foundAuthToken
		})

		It("updates the service auth token with the provided args", func() {
			runCommand("a label", "a provider", "a value")

			expectedAuthToken := models.ServiceAuthTokenFields{}
			expectedAuthToken.Guid = "found-auth-token-guid"
			expectedAuthToken.Label = "found label"
			expectedAuthToken.Provider = "found provider"
			expectedAuthToken.Token = "a value"

			Expect(ui.Outputs).To(ContainSubstrings(
				[]string{"Updating service auth token as", "my-user"},
				[]string{"OK"},
			))

			Expect(authTokenRepo.FindByLabelAndProviderLabel).To(Equal("a label"))
			Expect(authTokenRepo.FindByLabelAndProviderProvider).To(Equal("a provider"))
			Expect(authTokenRepo.UpdatedServiceAuthTokenFields).To(Equal(expectedAuthToken))
			Expect(authTokenRepo.UpdatedServiceAuthTokenFields).To(Equal(expectedAuthToken))
		})
	})
})
