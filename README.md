*Form3 exercise*
# Interview *accountapi*

 A client library in Go to access Form 3 fake account API service. Form 3 interview process exercise solution by Ruben Eguiluz.

## What's included

- The API client implementation, forlder client. 
- API resources schema, folder resources. Basically Form 3 API schema and builders for account resource.
- Unit and end2end test suites.
- End2End test execution infrastructure.

## Tests execution

Unit tests
``` make unitTest
```

* document technical decisions
  * ginkgo as BDD testing library because it's a good tool to write more readable test
  * gomega as matcher library, easy to use and fits well with ginkgo.
* Account resource: Create, Fetch, List and Delete
* contract test

TODO:

* Retry policy
* With context
* timeout?
* protobuffer
