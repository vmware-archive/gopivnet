package resource_test

import (
	"fmt"
	"net/http"

	"github.com/cfmobile/gopivnet/resource"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Releases", func() {
	var (
		req                resource.ReleaseRequester
		server             *ghttp.Server
		testRelease        *resource.Release
		testProductFiles   *resource.ProductFiles
		pivotalProductFile resource.ProductFile
		licenseProductFile resource.ProductFile
		eulaMessage        resource.EulaMessage
	)

	verifyHeaders := ghttp.CombineHandlers(
		ghttp.VerifyHeaderKV("Authorization", "Token token"),
		ghttp.VerifyHeaderKV("Content-Type", "application/json"),
		ghttp.VerifyHeaderKV("Accept", "application/json"),
		ghttp.VerifyHeaderKV("User-Agent", fmt.Sprintf("gopivnet %s", resource.Version)),
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		req = resource.NewRequester(server.URL(), "token")

		testRelease = &resource.Release{
			Id:      123,
			Version: "1.1",
			Links: resource.Links{
				"product_files": resource.Link{Url: server.URL() + "/api/v2/products/my-prod/releases/123/product_files"},
			},
		}

		pivotalProductFile = resource.ProductFile{
			Id:           200,
			AwsObjectKey: "my-prod.pivotal",
			FileVersion:  "2.0.0",
			Links: resource.Links{
				"download": resource.Link{Url: server.URL() + "/api/v2/products/my-prod/releases/123/product_files/200/download"},
			},
		}

		licenseProductFile = resource.ProductFile{
			Id:           201,
			AwsObjectKey: "my-prod.license",
			FileVersion:  "2.0.0",
			Links: resource.Links{
				"download": resource.Link{Url: server.URL() + "/api/v2/products/my-prod/releases/123/product_files/201/download"},
			},
		}

		testProductFiles = &resource.ProductFiles{
			Files: []resource.ProductFile{
				pivotalProductFile,
				licenseProductFile,
			},
		}

		eulaMessage = resource.EulaMessage{
			Status:  451,
			Message: "need to sign eula",
			Links: resource.Links{
				"eula_agreement": resource.Link{Url: server.URL() + "/api/v2/products/my-prod/releases/123/eula_acceptance"},
			},
		}
	})

	Context("GetProduct", func() {
		It("return an error if the token is not valid", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v2/products/my-prod/releases"),
					verifyHeaders,
					ghttp.RespondWith(http.StatusUnauthorized, ""),
				),
			)

			_, err := req.GetProduct("my-prod")
			Expect(err).To(HaveOccurred())

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		It("returns an error if the response is not valid", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v2/products/my-prod/releases"),
					verifyHeaders,
					ghttp.RespondWith(http.StatusOK, ""),
				),
			)

			_, err := req.GetProduct("my-prod")
			Expect(err).To(HaveOccurred())

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		It("returns the product returned by the server", func() {
			serverProd := resource.Product{
				Releases: []resource.Release{
					resource.Release{
						Id:      12,
						Version: "some-version",
					},
				},
			}

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v2/products/my-prod/releases"),
					verifyHeaders,
					ghttp.RespondWithJSONEncoded(http.StatusOK, serverProd),
				),
			)

			prod, err := req.GetProduct("my-prod")

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(err).ToNot(HaveOccurred())
			Expect(prod).To(Equal(&serverProd))
		})
	})

	Context("GetProductFiles", func() {
		It("returns an error if the release doesn't have product_files", func() {
			delete(testRelease.Links, "product_files")

			productFiles, err := req.GetProductFiles(*testRelease)
			Expect(server.ReceivedRequests()).To(HaveLen(0))

			Expect(productFiles).To(BeNil())
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the token is not valid", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v2/products/my-prod/releases/123/product_files"),
					verifyHeaders,
					ghttp.RespondWith(http.StatusUnauthorized, ""),
				),
			)

			productFiles, err := req.GetProductFiles(*testRelease)

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(productFiles).To(BeNil())
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if the result can't be parsed", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v2/products/my-prod/releases/123/product_files"),
					verifyHeaders,
					ghttp.RespondWith(http.StatusOK, ""),
				),
			)

			productFiles, err := req.GetProductFiles(*testRelease)

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(productFiles).To(BeNil())
			Expect(err).To(HaveOccurred())
		})

		It("returns the product files", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v2/products/my-prod/releases/123/product_files"),
					verifyHeaders,
					ghttp.RespondWithJSONEncoded(http.StatusOK, *testProductFiles),
				),
			)

			productFiles, err := req.GetProductFiles(*testRelease)

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(productFiles).To(Equal(testProductFiles))
			Expect(err).ToNot(HaveOccurred())
		})
	})

	Context("GetProductDownloadUrl", func() {
		It("returns an error if the token is not valid", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/api/v2/products/my-prod/releases/123/product_files/200/download"),
					verifyHeaders,
					ghttp.RespondWith(http.StatusUnauthorized, ""),
				),
			)

			url, err := req.GetProductDownloadUrl(&pivotalProductFile)
			Expect(err).To(HaveOccurred())

			Expect(server.ReceivedRequests()).To(HaveLen(1))
			Expect(url).To(Equal(""))
		})

		It("returns an error if the product file does not have a downlaod link", func() {
			delete(pivotalProductFile.Links, "download")

			url, err := req.GetProductDownloadUrl(&pivotalProductFile)
			Expect(err).To(HaveOccurred())

			Expect(server.ReceivedRequests()).To(HaveLen(0))
			Expect(url).To(Equal(""))
		})

		It("returns the download redirect url", func() {
			returnHeader := http.Header{}
			returnHeader.Add("Location", "testUrl")
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/api/v2/products/my-prod/releases/123/product_files/200/download"),
					verifyHeaders,
					ghttp.RespondWith(http.StatusFound, "", returnHeader),
				),
			)

			url, err := req.GetProductDownloadUrl(&pivotalProductFile)
			Expect(err).ToNot(HaveOccurred())
			Expect(server.ReceivedRequests()).To(HaveLen(1))

			Expect(url).To(Equal("testUrl"))
		})

		It("if the eula is not signed, it makes a request to the eula url", func() {
			returnHeader := http.Header{}
			returnHeader.Add("Location", "testUrl")

			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/api/v2/products/my-prod/releases/123/product_files/200/download"),
					verifyHeaders,
					ghttp.RespondWithJSONEncoded(resource.RequireEula, eulaMessage),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/api/v2/products/my-prod/releases/123/eula_acceptance"),
					verifyHeaders,
					ghttp.RespondWith(http.StatusOK, ""),
				),
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/api/v2/products/my-prod/releases/123/product_files/200/download"),
					verifyHeaders,
					ghttp.RespondWith(http.StatusFound, "", returnHeader),
				),
			)

			url, err := req.GetProductDownloadUrl(&pivotalProductFile)
			Expect(err).ToNot(HaveOccurred())

			Expect(server.ReceivedRequests()).To(HaveLen(3))
			Expect(url).To(Equal("testUrl"))
		})
	})
})
