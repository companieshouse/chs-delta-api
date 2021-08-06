# Unit testing a new schema

## Overview
As the chs-delta-api uses a 3rd party library (See `/docs` directory for documentation on the kin-openAPI3 library) to 
validate its requests, there is no way to directly unit test the code. This document will guide you through creating 
unit tests which assert that your open API schema is working correctly.

## 1. Where to add your unit tests
Move into the `/validation/schema_testing` directory and create a new directory to contain your unit tests. Inside of your
new directory, create a set `/request_bodies` and `/response_bodies` directories. These will contain your sample request 
and response bodies used to unit test the schema.

Finally, create a new Go file ending in `_test` (e.g. `exampleDeltaSchema_test.go`). This file will contain the unit 
tests used to test the schema.

## 2. Unit test structure
We elected to structure the open API schema unit tests using a `Karate API Scenario` layout. Each unit test will have a 
`Given`, `When`, `Then` stage and will use sample request bodies and response bodies to test the schema.

All unit tests will follow the same structure:

1. Load in the wanted request body from a text file using the `common.ReadRequestBody()` function
2. Create a test httpRequest and set its target to your delta's endpoint, and it's request body to the previously read in
request body, and finally set the requests headers.
3. Create an instance of the CHValidator and call it's `ValidateRequestAgainstOpenApiSpec` method, passing it your example
`request`, `open API spec location` and a dummy `contextId`.
4. The final step of every unit test is to read in your expected response body using the `common.ReadRequstBody()` function
and using the `common.CompareActualToExpected(actualResponseBody, expectedResponseBody)` to assert that your expected response
matches the actual response.

Example structure of a unit test
```go
func TestSchemaExample(t *testing.T) {

	Convey("Given I want to test the example-delta API schema", t, func() {

		exampleRequestBody := common.ReadRequestBody(exampleRequestBodyLocation)

		r := httptest.NewRequest("POST", exampleEndpoint, bytes.NewBuffer(exampleRequestBody))
		r = common.SetHeaders(r)

		Convey("When I call to validate the request body, providing a request", func() {

			chv := validation.NewCHValidator()

			validationErrs, _ := chv.ValidateRequestAgainstOpenApiSpec(r, apiSpecLocation, contextId)

			Convey("Then I am given a validation response", func() {
				exampleErrorResponseBody := common.ReadRequestBody(exampleErrorResponseBodyLocation)
				match := common.CompareActualToExpected(validationErrs, exampleErrorResponseBody)
                
				So(validationErrs, ShouldNotBeNil)
				So(match, ShouldEqual, true)
			})
		})
	})
}
```

## 3. What to cover and what not to cover in unit tests
The following areas of validation need to be covered by unit tests:

- Validate requests return no errors
- Mandatory / Required
- Type assertion (`int` only allows `int`, `string` only allows `string`)

The following areas of validation should not be covered by unit tests:

- Range validation (as it is better covered by Karate)