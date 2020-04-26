*Form3 exercise*
# Interview *accountapi*

 A client library in Go to access Form 3 fake account API service. Form 3 interview process exercise solution by Ruben Eguiluz.

## What's included

- The API client implementation, forlder client. 
- API resources schema, folder resources. Basically Form 3 API schema and builders for account resource.
- Unit and end2end test suites.
- End2End test execution infrastructure.

## Tests execution

Used libraries:

- Ginkgo: BDD testing library.
- Gomega: Matcher library.
- Gomock: Mocking library.
- Google uuid: Library used to generate random uuid values.

Before run tests, install testing library dependencies:
```
make deps 
```
Unit tests
```
make unitTest
```
End2end tests
```
make e2eTest
```

## Client use example

Create new Account resource example:

```go
    baseURL := "https://api.form3.tech//v1"
    client := NewForm3APIClient(baseURL, http.DefaultClient)
    
    id := "account-uuid"
    organisationID := "organisation-uuid"
    accountAttributes := map[string]interface{}{
		"country":       "GB",
		"base_currency": "GBP",
		"bank_id":       "400300",
		"bank_id_code":  "GBDSC",
		"bic":           "NWBKGB22",
	}
    accountData := resources.NewAccount(id, organisationID, accountAttributes)
    data := resources.NewDataContainer(accountData)

		resp, err := client.Create(
			context.Background(),
      resources.Account,
      data,
		)
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
