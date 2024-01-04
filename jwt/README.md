<h1>github.com/smxlong/kit/jwt</h1>

This package provides JWT functionality for Go's `net/http` package.

This package provides a JWT (JSON Web Token) middleware for Go's `net/http` package. It extracts and validates a JWT from the Authorization header of an HTTP request.

- [JWT Middleware](#jwt-middleware)
  - [Features](#features)
  - [Installation](#installation)
  - [Usage](#usage)
    - [Importing the package](#importing-the-package)
    - [Creating a new middleware instance](#creating-a-new-middleware-instance)
    - [Keys](#keys)
    - [Using custom claims](#using-custom-claims)
    - [Wrapping your HTTP handlers with the middleware](#wrapping-your-http-handlers-with-the-middleware)
    - [Accessing the claims in your handler](#accessing-the-claims-in-your-handler)
  - [Dependencies](#dependencies)
    - [`golang-jwt/jwt`](#golang-jwtjwt)

# JWT Middleware

## Features

- Extracts JWT from the Authorization header
- Validates JWT using a provided public key
- Supports custom claims
- Allows setting expected audience and issuer
- Stores parsed claims in the request context

## Installation

To install this package, you can use `go get`:

```bash
go get github.com/smxlong/kit/jwt
```

## Usage

### Importing the package

First, import the package:

```go
import "github.com/smxlong/kit/jwt"
```

You'll also need to import `github.com/golang-jwt/jwt/v5` under a different name
in order to access the claims:

```go
import gojwt "github.com/golang-jwt/jwt/v5"
```

### Creating a new middleware instance

The simplest configuration is to provide a key to use for JWT validation and
nothing else:

```go
jwtMiddleware := jwt.NewMiddleware(
    jwt.WithKey(yourKey),
)
```

Aside from validating the JWT signature, you can also set the expected audience
and issuer (which is recommended):

```go
jwtMiddleware := jwt.NewMiddleware(
    jwt.WithKey(yourKey),
    jwt.WithAudience(yourAudience),
    jwt.WithIssuer(yourIssuer),
)
```

### Keys

The key can be a byte slice (for symmetric signature algorithms), a keyfunc, or
a no-arg function that returns the key.

**Byte slice**

This is useful for HS256, HS384, and HS512 algorithms.

```go
jwtMiddleware := jwt.NewMiddleware(
    jwt.WithKey([]byte("your-key")),
)
```

**RSA public key**

This is useful for RS256, RS384, and RS512 algorithms.

```go
jwtMiddleware := jwt.NewMiddleware(
    jwt.WithKey(yourPublicKey),
)
```

**Keyfunc**

This is useful when you have multiple keys, and you want to select the key based
on the token's claims. For example, you might have a key per user, and you want
to select the key based on the user ID in the token's claims.

```go
jwtMiddleware := jwt.NewMiddleware(
    jwt.WithKey(func(token *gojwt.Token) (interface{}, error) {
        if token.Claims.(gojwt.MapClaims)["kid"] == "your-key-id" {
            return []byte("your-key"), nil
        }
    }),
)
```

**No-arg function**

This is useful when the keyfunc doesn't need to access the token's claims, but
it does need to be dynamic in some way. For example, loading the key from a file
each time it's needed.

```go
jwtMiddleware := jwt.NewMiddleware(
    jwt.WithKey(func() (interface{}, error) {
        return os.ReadFile("your-key-file")
    }),
)
```


### Using custom claims

By default, the middleware expects only the standard ("registered") claims as
defined in the
[`RegisteredClaims`](https://github.com/golang-jwt/jwt/blob/main/registered_claims.go)
struct of the `golang-jwt/jwt` package. If you want to use custom claims, you
can do so by implementing your own claims struct and a function that returns a
new instance of it, then passing that function to `NewMiddleware` using the
`WithNewClaims` option:

```go
type CustomClaims struct {
    gojwt.RegisteredClaims // always embed RegisteredClaims
    CustomClaim string `json:"custom_claim"` // add your custom claims
}

func newCustomClaims() gojwt.Claims {
    return &CustomClaims{}
}

jwtMiddleware := jwt.NewMiddleware(
    jwt.WithNewClaims(newCustomClaims),
    // ...
)
```

### Wrapping your HTTP handlers with the middleware

Call `Wrap` to wrap your HTTP handlers with the middleware:

```go
http.Handle("/path", jwtMiddleware.Wrap(yourHandler))
```

### Accessing the claims in your handler

In your handler, you can access the claims like this:

```go
claims, ok := r.Context().Value(jwt.ContextKeyClaims).(gojwt.Claims)
if !ok {
    // no claims found
}
```

Or, if you're using custom claims:

```go
claims, ok := r.Context().Value(jwt.ContextKeyClaims).(*CustomClaims)
if !ok {
    // no claims found, or claims are not of type *CustomClaims
}
```

## Dependencies

### `golang-jwt/jwt`

This package depends on the
[`github.com/golang-jwt/jwt/v5`](https://github.com/golang-jwt/jwt) package.
You'll need to make use of its
[`Claims`](https://github.com/golang-jwt/jwt/blob/main/claims.go) interface and
[`RegisteredClaims`](https://github.com/golang-jwt/jwt/blob/main/registered_claims.go)
struct, and optionally its
[`Keyfunc`](https://github.com/golang-jwt/jwt/blob/main/token.go) type, in order
to implement your own claims and key functions.
