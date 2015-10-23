package api_test

import (
	"errors"
	"io/ioutil"
	"net/http"
	"os"

	pivnetapi "github.com/cfmobile/gopivnet/api"
	"github.com/cfmobile/gopivnet/resource"
	"github.com/cfmobile/gopivnet/resource/fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Api", func() {
	var api *pivnetapi.PivnetApi
	var requester *fakes.FakeReleaseRequester
	var prod *resource.Product
	var productFiles *resource.ProductFiles

	BeforeEach(func() {
		prod = &resource.Product{
			Releases: []resource.Release{
				resource.Release{
					Id:      2,
					Version: "2.0",
				},
				resource.Release{
					Id:      1,
					Version: "1.0",
				},
			},
		}

		productFiles = &resource.ProductFiles{
			Files: []resource.ProductFile{
				resource.ProductFile{
					Id:           21,
					AwsObjectKey: "readme",
				},
				resource.ProductFile{
					Id:           22,
					AwsObjectKey: "product.pivotal",
				},
			},
		}

		requester = new(fakes.FakeReleaseRequester)
		requester.GetProductReturns(prod, nil)
		requester.GetProductFilesReturns(productFiles, nil)

		api = &pivnetapi.PivnetApi{
			Requester: requester,
		}

	})

	Context("GetLatestProductFile", func() {
		It("returns an error if there is no product name", func() {
			res, err := api.GetLatestProductFile("")

			Expect(res).To(BeNil())
			Expect(err).To(HaveOccurred())
		})

		It("fetches the right product", func() {
			api.GetLatestProductFile("myprod")

			Expect(requester.GetProductCallCount()).To(Equal(1))
			Expect(requester.GetProductArgsForCall(0)).To(Equal("myprod"))
		})

		It("returns an error if fetching a product fails", func() {
			requester.GetProductReturns(nil, errors.New("err"))
			res, err := api.GetLatestProductFile("myprod")

			Expect(res).To(BeNil())
			Expect(err).To(HaveOccurred())
		})

		It("fetches the product files from the returned product", func() {
			api.GetLatestProductFile("myprod")

			Expect(requester.GetProductFilesCallCount()).To(Equal(1))
			Expect(requester.GetProductFilesArgsForCall(0)).To(Equal(prod.Releases[0]))
		})

		It("returns an error if GetProductFiles fails", func() {
			requester.GetProductFilesReturns(nil, errors.New("err"))
			res, err := api.GetLatestProductFile("myprod")

			Expect(res).To(BeNil())
			Expect(err).To(HaveOccurred())
		})

		It("returns the latest product file", func() {
			res, err := api.GetLatestProductFile("myprod")

			Expect(res).To(Equal(&productFiles.Files[1]))
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error if there's no pivotal product", func() {
			productFiles.Files = productFiles.Files[:1]
			requester.GetProductFilesReturns(productFiles, nil)

			res, err := api.GetLatestProductFile("myprod")

			Expect(res).To(BeNil())
			Expect(err).To(HaveOccurred())
		})
	})

	Context("GetProductFileForVersion", func() {
		It("returns an error if there is no product name", func() {
			res, err := api.GetProductFileForVersion("", "1.0")

			Expect(res).To(BeNil())
			Expect(err).To(HaveOccurred())
		})

		It("returns an error if there is no product version", func() {
			res, err := api.GetProductFileForVersion("name", "")

			Expect(res).To(BeNil())
			Expect(err).To(HaveOccurred())
		})

		It("fetches the right product", func() {
			api.GetProductFileForVersion("myprod", "1.0")

			Expect(requester.GetProductCallCount()).To(Equal(1))
			Expect(requester.GetProductArgsForCall(0)).To(Equal("myprod"))
		})

		It("returns an error if fetching a product fails", func() {
			requester.GetProductReturns(nil, errors.New("err"))
			res, err := api.GetProductFileForVersion("myprod", "1.0")

			Expect(res).To(BeNil())
			Expect(err).To(HaveOccurred())
		})

		It("fetches the product files from the returned product", func() {
			api.GetProductFileForVersion("myprod", "1.0")

			Expect(requester.GetProductFilesCallCount()).To(Equal(1))
			Expect(requester.GetProductFilesArgsForCall(0)).To(Equal(prod.Releases[1]))
		})

		It("returns an error if GetProductFiles fails", func() {
			requester.GetProductFilesReturns(nil, errors.New("err"))
			res, err := api.GetProductFileForVersion("myprod", "1.0")

			Expect(res).To(BeNil())
			Expect(err).To(HaveOccurred())
		})

		It("returns the latest product file", func() {
			res, err := api.GetProductFileForVersion("myprod", "1.0")

			Expect(res).To(Equal(&productFiles.Files[1]))
			Expect(err).ToNot(HaveOccurred())
		})

		It("returns an error if there's no pivotal product", func() {
			productFiles.Files = productFiles.Files[:1]
			requester.GetProductFilesReturns(productFiles, nil)

			res, err := api.GetProductFileForVersion("myprod", "1.0")

			Expect(res).To(BeNil())
			Expect(err).To(HaveOccurred())
		})
	})

	Context("Download", func() {
		var file *os.File
		var server *ghttp.Server

		BeforeEach(func() {
			var err error
			file, err = ioutil.TempFile("", "")
			file.Close()
			Expect(err).ToNot(HaveOccurred())

			server = ghttp.NewServer()
			server.AppendHandlers(
				ghttp.RespondWith(http.StatusOK, `aaa`),
			)
		})

		AfterEach(func() {
			os.Remove(file.Name())
			server.Close()
		})

		var testFileIsEmpty = func() {
			fileInfo, _ := os.Lstat(file.Name())
			ExpectWithOffset(1, fileInfo.Size()).To(Equal(int64(0)))
		}

		It("returns an error if the product is nil", func() {
			err := api.Download(nil, file.Name())
			Expect(err).To(HaveOccurred())
			testFileIsEmpty()
		})

		It("returns an error if it can't get the product download url", func() {
			requester.GetProductDownloadUrlReturns("", errors.New("err"))
			err := api.Download(nil, file.Name())
			Expect(err).To(HaveOccurred())
			testFileIsEmpty()
		})

		It("downloads the data at the url", func() {
			requester.GetProductDownloadUrlReturns(server.URL(), nil)
			err := api.Download(&resource.ProductFile{}, file.Name())

			Expect(err).ToNot(HaveOccurred())
			res, err := ioutil.ReadFile(file.Name())
			Expect(res).To(Equal([]byte("aaa")))
		})
	})
})
