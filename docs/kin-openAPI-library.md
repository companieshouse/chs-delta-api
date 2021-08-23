# Kin openAPI3 validation library documentation

## Overview
this document will provide useful information about the kin-openAPI3 validation library that was used in the chs-delta-api.

## Errors
The kin-openAPI3 library uses Go version 1.14+ to wrap the standard Go Error struct to allow for more detail to be returned 
to the user. This means that the caller of the kin-openAPI3 validator will need to parse the information returned out of 
the wrapped error(s) to get the necessary values back.

The chs-delta-api takes the wrapped error(s) returned from the kin-openAPI3 validator and parses them into the standard 
[CH Errors Object](https://developer-specs.company-information.service.gov.uk/companies-house-public-data-api/resources/error?v=latest) by using the errors.Is and errors.As functions introduced to Go in version 1.13. To view 
the code implementation of this, please see the `/validation/chOpenApiValidator.go` file; specfically the `getCHErrors` 
function.

## Supported specs
For information on the spec types supported by the kin-openAPI3 library, please visit the following URL: 
https://github.com/getkin/kin-openapi#structure.

