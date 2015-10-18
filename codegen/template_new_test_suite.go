package codegen

// GetNewTestSuiteTemplate => builds the template string for generating the web server test suite file
func GetNewTestSuiteTemplate(cc CodeConfig) string {
	return `
	package webserver

  {{.TimeStampHeader false}}

	import (
		"bytes"
		"encoding/json"
		"log"
		"net/http"
		"net/http/httptest"
		"reflect"
		"sort"
		"strconv"
		"strings"
		"testing"

		"github.com/obieq/gas"
		gapi "github.com/obieq/gson-api"
		. "github.com/onsi/ginkgo"
		. "github.com/onsi/gomega"
		"github.com/pallinder/go-randomdata"
		"gopkg.in/guregu/null.v3"
	)

	const HTTP_POST  = "POST"
  const HTTP_PATCH = "PATCH"

	type Attributes struct {
	  Attributes interface{} ` + "`" + `json:"attributes"` + "`" + `
  }

  type JSONRequest struct {
	  Data interface{} ` + "`" + `json:"data"` + "`" + `
  }

	var (
		response  *httptest.ResponseRecorder
	)

	func TestMain(t *testing.T) {
		RegisterFailHandler(Fail)
		RunSpecs(t, "Main Suite")
	}

	var _ = BeforeSuite(func() {
	})

	var _ = AfterSuite(func() {
  })

	func parseJSON(recorder *httptest.ResponseRecorder, key string, isAttribute bool) interface{} {
		j := gapi.JsonApiResource{}
		json.Unmarshal(recorder.Body.Bytes(), &j)
		data := j.Data.(map[string]interface{})
		attributes := data["attributes"].(map[string]interface{})

		if isAttribute {
			return attributes[key]
		}

		return data[key]
	}

	func parseCreatedAt(recorder *httptest.ResponseRecorder) string {
		// NOTE: due to timestamp differences (long story), we need to replace the timestamp
		// returned by POST with the timestamp returned by GET
		return strings.Split(strings.Split(recorder.Body.String(), ` + "`" + `"created-at":"` + "`" + `)[1], ` + "`" + `"` + "`" + `)[0]
	}

	func expectedEmptyGetResponse() string {
	  return ` + "`" + `{"data":[]}` + "`" + `
  }

	func expectedNotFoundResponse() string {
	  return ` + "`" + `{ "errors": { "status": "404", "detail": "record not found" } }` + "`" + `
  }

  func expectedGetOneResponse(resourceName string, orderedKeys map[int]string, attrs map[string]interface{}) string {
    // parse id
    id := attrs["id"].(string)
    delete(attrs, "id") // NOTE: very important step!

	  return ` + "`" + `{` + "`" + ` +
		  ` + "`" + `"data":{"type":"` + "`" + " + resourceName + " + "`" + `","id":"` + "`" + ` + id + ` + "`" + `",` + "`" + ` +
		  ` + "`" + `"attributes":{` + "`" + ` + parseKVs(orderedKeys, attrs) + ` + "`" + `,"updated-at":null}}}` + "`" + `
  }

func expectedGetAllResponse(resourceName string, orderedKeys map[int]string, attrs map[string]interface{}) string {
  // parse id
  id := attrs["id"].(string)
  delete(attrs, "id") // NOTE: very important step!

	return ` + "`" + `{` + "`" + ` +
		` + "`" + `"data":[{"type":"` + "`" + " + resourceName + " + "`" + `","id":"` + "`" + ` + id + ` + "`" + `",` + "`" + ` +
		` + "`" + `"attributes":{` + "`" + ` + parseKVs(orderedKeys, attrs) + ` + "`" + `,"updated-at":null}}]}` + "`" + `
}

func expectedPostPatchResponse(resourceName string, orderedKeys map[int]string, attrs map[string]interface{}) string {
  // parse id
  id := attrs["id"].(string)
  delete(attrs, "id") // NOTE: very important step!

  // parse updated at
  updatedAt := ""
  if attrs["updated-at"] == nil {
    updatedAt = ` + "`" + `,"updated-at":null` + "`" + `
  }

	return "{" +
		` + "`" + `"data":{"type":"` + "`" + " + resourceName + " + "`" + `","id":"` + "`" + ` + id + ` + "`" + `",` + "`" + ` +
		` + "`" + `"attributes":{` + "`" + ` + parseKVs(orderedKeys, attrs) + updatedAt + ` + "`" + `}}}` + "`" + `
}

func parseKVs(orderedKeys map[int]string, attrs map[string]interface{}) string {
	s := ""
	sorted := make([]int, 0, len(orderedKeys))
	for i := range orderedKeys {
		sorted = append(sorted, i)
	}
	sort.Ints(sorted)

	for _, i := range sorted {
		k := orderedKeys[i]
		v := attrs[k]
		s += ` + "`" + `"` + "`" + ` + gas.String(k).Dasherize() + ` + "`" + `":` + "`" + `
		switch v.(type) {
		case string:
			s += ` + "`" + `"` + "`" + ` + v.(string) + ` + "`" + `"` + "`" + `
		case int:
			s += strconv.Itoa(v.(int))
		case bool:
			s += strconv.FormatBool(v.(bool))
		case int64:
			s += strconv.FormatInt(v.(int64), 10)
		case float32:
			s += strconv.FormatFloat(float64(v.(float32)), 'f', -1, 32)
		case float64:
			s += strconv.FormatFloat(v.(float64), 'f', -1, 64)
		default:
			log.Panicln("No switch for KV type:", reflect.TypeOf(v))
		}
		if i < len(attrs)-1 {
			s += ","
		}
	}
	return s
}

func makeRESTRequest(verb string, routeName string, resource gapi.Resourcer) (*httptest.ResponseRecorder, map[string]interface{}) {
  kvs := map[string]interface{}{}
  generateRandomData(kvs, resource)

	// set attributes struct values w/ the random data values
  val := reflect.Indirect(reflect.ValueOf(resource))
	for i := 0; i < val.Type().NumField(); i++ {
		f := val.Type().Field(i)

		// skip both created at and updated at
		if f.Name == "CreatedAt" || f.Name == "UpdatedAt" || f.Name == "Resource" {
      continue
		}

		switch f.Type.String() {
		case "null.String":
			v := kvs[gas.String(f.Name).Dasherize()].(string)
			val.FieldByName(f.Name).Set(reflect.ValueOf(null.StringFrom(v)))
    case "null.Bool":
			v := kvs[gas.String(f.Name).Dasherize()].(bool)
		  val.FieldByName(f.Name).Set(reflect.ValueOf(null.BoolFrom(v)))
		case "null.Int":
			v := kvs[gas.String(f.Name).Dasherize()].(int)
			val.FieldByName(f.Name).Set(reflect.ValueOf(null.IntFrom(int64(v))))
		default:
			log.Panicln("reflect.Set not mapped for:", f.Type.String())
		}
	}

  return sendTestRequest(verb, routeName, resource, kvs)
}

func sendTestRequest(
	httpVerb string,
	routeName string,
	resource gapi.Resourcer,
	kvs map[string]interface{}) (*httptest.ResponseRecorder, map[string]interface{}) {

	a := Attributes{Attributes: resource}
  j := JSONRequest{Data: a}

	body, err := json.Marshal(j)
	if err != nil {
		log.Println("Unable to marshal Resource:", resource)
	}

	server := NewServer()
	request, _ := http.NewRequest(httpVerb, routeName, bytes.NewReader(body))
	recorder := httptest.NewRecorder()
	server.ServeHTTP(recorder, request)

	// get id and created-at values
	kvs["id"] = parseJSON(recorder, "id", false).(string)
  kvs["created-at"] = parseJSON(recorder, "created-at", true).(string)
	if httpVerb == HTTP_PATCH {
		kvs["updated-at"] = parseJSON(recorder, "updated-at", true).(string)
	}

	return recorder, kvs
}

func generateRandomData(kvs map[string]interface{}, resource gapi.Resourcer) {
	val := reflect.Indirect(reflect.ValueOf(resource))
	for i := 0; i < val.Type().NumField(); i++ {
		f := val.Type().Field(i)

		// skip both created at and updated at
		if f.Name == "CreatedAt" || f.Name == "UpdatedAt" || f.Name == "Resource" {
      continue
		}

		switch f.Type.String() {
		case "string", "*string", "null.String":
			kvs[gas.String(f.Name).Dasherize()] = randomdata.SillyName()
		case "int", "int64", "*int", "*int64", "null.Int":
			kvs[gas.String(f.Name).Dasherize()] = randomdata.Number(20000)
		case "bool", "*bool", "null.Bool":
			kvs[gas.String(f.Name).Dasherize()] = randomdata.Boolean()
		case "float32", "float64", "*float32", "*float64", "null.Float":
			kvs[gas.String(f.Name).Dasherize()] = randomdata.Decimal(20000)
		default:
			log.Panicln("No random data mapping for:", f.Name, f.Type)
		}
	}
}`
}
