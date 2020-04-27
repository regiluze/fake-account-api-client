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
End2end tests (before running this job the docker-compose should be running)
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

    resp, err := client.Create(context.Background(), resources.Account, data)
```
## Technical decisions

- Ginkgo as BDD testing library because it's a good tool to write more readable tests.
- Gomega as matcher library, because it's easy to read and fits well with Ginkgo.
- The client implementation is done by doing TDD and to isolite the client, I created HTTPClient interface, with just Do(req *http.Request) method. In this way, I could mock the http client and start the iterations.
- There are more unit tests than e2e tests, the only reason for that it's because doing TDD you write a lot of tests to cover completelly the system under test and not because they are more important than e2e tests. In my opinion the good ones are the e2e tests because it's the real interaction with the API and they are a kind of contract tests and they are less coupled to the implementation.
- There are 3 custom error types:
  - ErrBadRequest: Server return status code is 400, I created this error type to be accesible the server return information about the problem.
  - ErrNotFound: Server return status code is 404, I created this error type to be easier to catch the error type when using the client. Sometimes when fetching a resource, the behaviour after getting this error it's different than getting another status code like 500 or 400, for instance if one you just need to check if the account exists or not the execution path would be different for 404 error than 500.
  - ErrResponseStatusCode: Server return status code is 40X (less 400 and 404) or 50X. The status code is accesible.
- There aren't any validation in the client, it's rely on server validation, in my opinion doesn't make sense to do the business validation in the client when the business knowledge is in the server and the business decisions are made in the server. For 'country' required account parameter, it returns an ErrBadRequest error with the information about the required parameter, there is a specific end2end test for this.
- Context parameter: At the begining my idea was to duplicate the client public API like CreateWithContext, and so on. But finally I decided to include the context as a parameter in all public methods because in my opinion the context in http request is a good practice because for instance you could include a timeout, some data for traceability, etc.
- I created Client interface type because it's useful when using it, for instance to be clear the contract of the client or to create a mock of the client.
- Public method names: At the begining the methods names were more coupled to the Account resource but during the implementation I realized that the Form3 API schema was generic so I decided to change them to more generic way, so the Account name went from the method name to as a parameter. In case of extending the client and support another resource, the changes would be just the resource and the endpoint mapping.
- I extracted the URL builder struct to keep all the code about the resources endpoints in one place.

## TODO list

- The rest of the resource methods and List filter feature.
- Timeout: Add init client method with timeout parameter to create a http client setting this parameter.
- Retry logic.
- Authentication.
- Use env variables to set the base URL value.
