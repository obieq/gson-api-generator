package codegen

// GetRequestTestTemplate => builds the template string for generating the request test file for a new resource
func GetRequestTestTemplate(cc CodeConfig) string {
	return `
package webserver

{{.TimeStampHeader false}}

import (
	"net/http"
	"net/http/httptest"

	"github.com/obieq/goar"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"{{.GoPath}}/app/models"
	"{{.GoPath}}/app/resources"
)

var _ = Describe("{{.EntityName}}", func() {
	var (
		db{{.EntityName}}        *models.{{.EntityName}} = goar.ToAR(&models.{{.EntityName}}{}).(*models.{{.EntityName}})
		g{{.EntityName}}Resource *resources.{{.EntityName}}
		server               Server
		request              *http.Request
		gRecorder            *httptest.ResponseRecorder
		gKVS                 = map[string]interface{}{}
	)

	BeforeEach(func() {
		var recorder *httptest.ResponseRecorder

		// instantiate a new web server
		server = NewServer()

		// Record HTTP responses
		gRecorder = httptest.NewRecorder()

		// Truncate table
		db{{.EntityName}}.Truncate() // delete all records created during previous test

		// Save a new {{.EntityName}} before each test
		g{{.EntityName}}Resource = &resources.{{.EntityName}}{}
		recorder, gKVS = makeRESTRequest(HTTP_POST, "/v1/{{.DasherizedEntityName true}}", g{{.EntityName}}Resource)
		Ω(recorder.Code).Should(Equal(201))
	})

	Context("GET", func() {
		Context("One", func() {
			BeforeEach(func() {
				// Set up a new GET request before every test
				request, _ = http.NewRequest("GET", "/v1/{{.DasherizedEntityName true}}/"+gKVS["id"].(string), nil)
			})

			Context("when no {{.EntityName}}s exist", func() {
				It("should return a 401 status code and an errors collection", func() {
					db{{.EntityName}}.Truncate() // delete all records created during previous test

					server.ServeHTTP(gRecorder, request)
					Ω(gRecorder.Code).Should(Equal(404))
					Ω(gRecorder.Body).Should(MatchJSON(expectedNotFoundResponse()))
				})
			})

			Context("when a {{.EntityName}} exists", func() {
				It("should return a 200 status code and a single json api resource", func() {
					server.ServeHTTP(gRecorder, request)
					Ω(gRecorder.Code).Should(Equal(200))

					createdAt := parseCreatedAt(gRecorder)

					// verify response
					orderedKeys := {{.ParseOrderedKeys true false}}
					gKVS["created-at"] = createdAt
					Ω(gRecorder.Body).Should(MatchJSON(expectedGetOneResponse("{{.DasherizedEntityName true}}", orderedKeys, gKVS)))
				})
			})
		})
		Context("All", func() {
			BeforeEach(func() {
				// Set up a new GET request before every test
				request, _ = http.NewRequest("GET", "/v1/{{.DasherizedEntityName true}}", nil)
			})

			Context("when no {{.EntityName}}s exist", func() {
				It("should return a 200 status code and an empty array", func() {
					db{{.EntityName}}.Truncate() // delete all records created during previous test

					server.ServeHTTP(gRecorder, request)
					Ω(gRecorder.Code).Should(Equal(200))
					Ω(gRecorder.Body).Should(MatchJSON(expectedEmptyGetResponse()))
				})
			})

			Context("when a {{.EntityName}} exists", func() {
				It("should return a 200 status code and an array", func() {
					server.ServeHTTP(gRecorder, request)
					Ω(gRecorder.Code).Should(Equal(200))

					createdAt := parseCreatedAt(gRecorder)

					// verify response
					orderedKeys := {{.ParseOrderedKeys true false}}
					gKVS["created-at"] = createdAt
					Ω(gRecorder.Body).Should(MatchJSON(expectedGetAllResponse("{{.DasherizedEntityName true}}", orderedKeys, gKVS)))
				})
			})
		})
	})

	Context("POST", func() {
		It("should return a 201 Status Code", func() {
			db{{.EntityName}}.Truncate() // delete all records created during previous test

			r, kvs := makeRESTRequest(HTTP_POST, "/v1/{{.DasherizedEntityName true}}", &resources.{{.EntityName}}{})
			Ω(r.Code).Should(Equal(201))

			// verify response
			orderedKeys := {{.ParseOrderedKeys true false}}
			Ω(r.Body).Should(MatchJSON(expectedPostPatchResponse("{{.DasherizedEntityName true}}", orderedKeys, kvs)))
		})
	})

	Context("PATCH", func() {
		It("should return a 200 status code", func() {
			resource := &resources.{{.EntityName}}{}
			r, kvs := makeRESTRequest(HTTP_PATCH, "/v1/{{.DasherizedEntityName true}}/"+gKVS["id"].(string), resource)

			// verify response
			orderedKeys := {{.ParseOrderedKeys true true}}
			Ω(r.Body).Should(MatchJSON(expectedPostPatchResponse("{{.DasherizedEntityName true}}", orderedKeys, kvs)))
		})
	})
})`
}
