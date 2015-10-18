package resource_test

import (
	. "github.com/cfmobile/gopivnet/resource"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Resource", func() {
	Context("ProductFile", func() {
		It("returns the product name", func() {
			productFile := ProductFile{
				Id:           1,
				AwsObjectKey: "abc/asd/wdv/test",
			}

			Expect(productFile.Name()).To(Equal("test"))
		})
	})
})
