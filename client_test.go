package accountclient

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	//. "github.com/onsi/gomega"
)

var _ = Describe("Account api resource client", func() {
	var (
		client *Form3Client
		//httpClientMock HTTPClientMock
	)

	BeforeEach(func() {
		client = NewForm3Client()
	})

	Describe("Resource operations", func() {
		Context("Create", func() {
			It("blablabla", func() {
				fmt.Println(">>>> client ", client)
			})
		})
	})
})
