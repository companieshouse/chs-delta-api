# Unit testing a new schema

## Overview
As the chs-delta-api uses a 3rd party library (See `/docs` directory for documentation on the kin-openAPI3 library) to 
validate its requests, there is no way to directly unit test the code. This document will guide you through creating 
unit tests which assert that your openAPI schema is working correctly.

## 1. Where to add your unit tests
Move into the `/validation/schema_testing` directory and create a new directory to contain your unit tests. Inside of your
new directory, create a set `/request_bodies` and `/response_bodies` directories. These will contain your sample request 
and response bodies used to unit test the schema.

Finally, create a new Go file ending in `_test` (e.g. `exampleDeltaSchema_test.go`). This file will contain the unit 
tests used to test the schema.

## 2. Unit test structure
We elected to structure the openAPI schema unit tests using a `Karate API Scenario` layout. Each unit test will have a 
`Given`, `When`, `Then` stage and will use sample request bodies and response bodies to test the schema.

All unit tests will follow the same structure:

1. Load in the example request body from a text file using the `common.ReadRequestBody()` function
2. Create and populate a test httpRequest 
    1. Set its target to your delta's endpoint
    2. Set it's body to that previously loaded in step 1
    3. Set the request's headers using the provided `SetHeaders` function in `/validation/schema_testing/common/schemaUnitTesting.go`
3. Create an instance of the CHValidator and call its `ValidateRequestAgainstOpenApiSpec` method, passing it your example
`request`, `openAPI spec location` and a dummy `contextId`.
4. The final step of every unit test is to read in your expected response body using the `common.ReadRequstBody()` function
and using the `common.CompareActualToExpected(actualResponseBody, expectedResponseBody)` to assert that your expected response
matches the actual response.

Example structure of a unit test (Mandatory elements missing unit test)
```go
func TestSchemaExampleMissingMandatory(t *testing.T) {

	Convey("Given I want to test that missing mandatory fields in the exampleDelta are correctly reported", t, func() {

		mandatoryMissingRequestBody := common.ReadRequestBody(mandatoryMissingRequestBodyLocation)

		r := httptest.NewRequest("POST", exampleEndpoint, bytes.NewBuffer(mandatoryMissingRequestBody))
		r = common.SetHeaders(r)

		Convey("When I call to validate the request body, providing a request which is missing mandatory fields", func() {

			chv := validation.NewCHValidator()

			actualResponseBody, _ := chv.ValidateRequestAgainstOpenApiSpec(r, apiSpecLocation, contextId)

			Convey("Then I am given a validation response stating that mandatory fields are missing", func() {
				expectedResponseBody := common.ReadRequestBody(mandatoryMissingErrorResponseBodyLocation)
				match := common.CompareActualToExpected(actualResponseBody, expectedResponseBody)
                
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