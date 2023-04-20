# About me

Hi, I'm Alan Gatticelli 😁.

I've been playing around with programming in the professional industry for 10 years. My main background is software development, but I also have experience in event-driven architectures in AWS, mainly serverless.

My current tech stack is Typescript + CDK but I also use Python for scripting.

Regarding my Golang experience, I met this amazing language 3 years ago and I would love to dive deep into it.

# Installation

In order to install the SDK, you need to run.

```bash
go get github.com/agatticelli/form3-client-go/form3
```

# Usage

First, let's import the library.

```go
import "github.com/agatticelli/form3-client-go/form3"
```

## Fetch an account

```go
client := form3.NewClient(nil)
account, _, _ := client.Account.Fetch(context.Background(), accountID)
```

## Create an account

```go
client := form3.NewClient(nil)
attributes := form3.CreateAccountAttributes{
  BankID:      "20041",
  BankIDCode:  "FR",
  Bic:         "NWBKFR42",
  Name:        []string{"Alan Gatticelli"},
  Country:     "FR",
}
account, _, err := client.Account.Create(
  context.Background(), accountID, organisationID, &attributes
)
```

## Delete an account

```go
client := form3.NewClient(nil)
err := client.Account.Delete(context.Background(), accountID, version)
```

# Contributing

In order to run all available tests, unit and integration, you need to be in the root path and start all the services with

```bash
docker compose up -d
```

After all services finished bootstraping, you should run this every time you want to run the tests.

```bash
docker compose up form3-client
```

Also, you can run unit tests only like this

```bash
cd form3
go test -v ./...
```

# Extras (for the Form3 team)

There are some features that were not developed as the challenge told so, but I will list some of them that I consider a production ready SDK should have.

## Rate Limit

It can happen that the Form3 API raises a **Rate Limit** error by returning a `429` status code. It would be great to have a setting to enable some sort of exponential backoff algorithm to automatically retry the request after a few seconds.

## Timeout

There is an example in the Client.Do() unit test for this, but it can be easily implemented in a real world scenario at the client or request level.
For example, at client level would be something like the following:

```go
httpClient := http.Client{Timeout: 2 * time.Second}
client := form3.NewClient(&httpClient)
account, _, _ := client.Account.Fetch(context.Background(), accountID)
```

Or it can be done at the request level doing something like:

```go
client := form3.NewClient(nil)
ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
defer cancel()
account, _, _ := client.Account.Fetch(ctx, accountID)
```

## Retries

Every HTTP client should have a retry strategy for unexpected errors. This errors are related to `5xx` and `429` status codes (if the service has a rate limit). So, our client can receive a `RetryCount` setting to enable it.

# How the solution was thought (for the Form3 team)

At first, I began with a simple `Client` struct with methods such as CreateAccount, FetchAccount and DeleteAccount. It also had some non-exported methods such as:

- `newRequest` to build a request from scratch with its method and url, encoding of body data, etc
- `decodeBody` to decode the responses into a valid struct or an error struct.
- `do` this function was in charge of creating the request with `newRequest`, dispatch it and bind the response by calling the `decodeBody` method.

After working a while with integration tests structure, I noticed that it would be great to have a way of purging all the previous created accounts to avoid pollution between each test. So I exported the `Do` method to allow users to implement API calls to non implemented endpoints. With that, I was able to call the `List Accounts` endpoint and delete one by one with the `DeleteAccount` method.

Finally, I realized that the Accounts API is one of the multiple services that Form3 offers, so I created an `AccountService` struct which receives a `form3.Client` instance to make the API calls and with that, I simplify the calls to `account.Create`, `acount.Fetch` and `account.Delete`.
