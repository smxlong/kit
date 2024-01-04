package jwt

import (
	"crypto/rand"
	"crypto/rsa"
	"net/http"
	"testing"

	gojwt "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func Test_That_NewMiddleware_Constructs_A_Middleware_Correctly_From_Options(t *testing.T) {
	var privateKey interface{} = "privateKey"
	var newClaimsCalled bool
	var newClaims func() gojwt.Claims = func() gojwt.Claims {
		newClaimsCalled = true
		return &gojwt.RegisteredClaims{}
	}
	var audience string
	var issuer string

	middleware := NewMiddleware(
		WithKey(privateKey),
		WithNewClaims(newClaims),
		WithAudience(audience),
		WithIssuer(issuer),
	)

	middleware.newClaims()

	assert.NotNil(t, middleware)
	assert.Equal(t, privateKey, middleware.publicKey)
	assert.True(t, newClaimsCalled)
	assert.Equal(t, audience, middleware.audience)
	assert.Equal(t, issuer, middleware.issuer)
}

func Test_That_getToken_Returns_An_Error_When_No_Token_Is_Present(t *testing.T) {
	middleware := NewMiddleware()

	token, err := middleware.getToken(&http.Request{})

	assert.Equal(t, "", token)
	assert.Equal(t, "no token", err.Error())
}

func Test_That_getToken_Returns_An_Error_When_The_Token_Is_Not_A_Bearer_Token(t *testing.T) {
	middleware := NewMiddleware()

	token, err := middleware.getToken(&http.Request{
		Header: http.Header{
			"Authorization": []string{"test"},
		},
	})

	assert.Equal(t, "", token)
	assert.Equal(t, "no token", err.Error())
}

func Test_That_getToken_Returns_The_Token_When_The_Token_Is_A_Bearer_Token(t *testing.T) {
	middleware := NewMiddleware()

	token, err := middleware.getToken(&http.Request{
		Header: http.Header{
			"Authorization": []string{"Bearer test"},
		},
	})

	assert.Equal(t, "test", token)
	assert.Nil(t, err)
}

func Test_That_parseClaims_Returns_An_Error_When_The_Token_Is_Empty(t *testing.T) {
	middleware := NewMiddleware()

	claims, err := middleware.parseClaims("")

	assert.Nil(t, claims)
	assert.Equal(t, "token is malformed: token contains an invalid number of segments", err.Error())
}

func Test_That_parseClaims_Returns_An_Error_When_The_Token_Is_Invalid(t *testing.T) {
	middleware := NewMiddleware()

	claims, err := middleware.parseClaims("test")

	assert.Nil(t, claims)
	assert.Equal(t, "token is malformed: token contains an invalid number of segments", err.Error())
}

func Test_That_parseClaims_Returns_The_Claims_When_The_Token_Is_Valid(t *testing.T) {
	middleware := NewMiddleware(
		WithKey([]byte("secret")),
	)

	token, err := gojwt.NewWithClaims(gojwt.SigningMethodHS256, &gojwt.RegisteredClaims{Subject: "test"}).SignedString([]byte("secret"))
	assert.NoError(t, err)

	claims, err := middleware.parseClaims(token)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "test", claims.(*gojwt.RegisteredClaims).Subject)
}

func Test_That_parseClaims_Returns_The_Claims_When_The_Token_Is_Valid_RSA(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	middleware := NewMiddleware(
		WithKey(&privateKey.PublicKey),
	)

	token, err := gojwt.NewWithClaims(gojwt.SigningMethodRS256, &gojwt.RegisteredClaims{Subject: "test"}).SignedString(privateKey)
	assert.NoError(t, err)

	claims, err := middleware.parseClaims(token)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "test", claims.(*gojwt.RegisteredClaims).Subject)
}

func Test_That_parseClaims_Returns_An_Error_When_The_Token_Is_Valid_But_The_Audience_Is_Wrong(t *testing.T) {
	middleware := NewMiddleware(
		WithKey([]byte("secret")),
		WithAudience("test"),
	)

	token, err := gojwt.NewWithClaims(gojwt.SigningMethodHS256, &gojwt.RegisteredClaims{Subject: "test"}).SignedString([]byte("secret"))
	assert.NoError(t, err)

	claims, err := middleware.parseClaims(token)

	assert.Nil(t, claims)
	assert.Equal(t, "token has invalid claims: token is missing required claim: aud claim is required", err.Error())
}

func Test_That_parseClaims_Returns_An_Error_When_The_Token_Is_Valid_But_The_Issuer_Is_Wrong(t *testing.T) {
	middleware := NewMiddleware(
		WithKey([]byte("secret")),
		WithIssuer("test"),
	)

	token, err := gojwt.NewWithClaims(gojwt.SigningMethodHS256, &gojwt.RegisteredClaims{Subject: "test"}).SignedString([]byte("secret"))
	assert.NoError(t, err)

	claims, err := middleware.parseClaims(token)

	assert.Nil(t, claims)
	assert.Equal(t, "token has invalid claims: token is missing required claim: iss claim is required", err.Error())
}

func Test_That_parseClaims_Returns_The_Claims_When_The_Token_Is_Valid_And_The_Audience_And_Issuer_Are_Correct(t *testing.T) {
	middleware := NewMiddleware(
		WithKey([]byte("secret")),
		WithAudience("test"),
		WithIssuer("test"),
	)

	token, err := gojwt.NewWithClaims(gojwt.SigningMethodHS256, &gojwt.RegisteredClaims{Subject: "test", Audience: []string{"test"}, Issuer: "test"}).SignedString([]byte("secret"))
	assert.NoError(t, err)

	claims, err := middleware.parseClaims(token)

	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, "test", claims.(*gojwt.RegisteredClaims).Subject)
	assert.Equal(t, 1, len(claims.(*gojwt.RegisteredClaims).Audience))
	assert.Equal(t, "test", claims.(*gojwt.RegisteredClaims).Audience[0])
	assert.Equal(t, "test", claims.(*gojwt.RegisteredClaims).Issuer)
}

func Test_That_Wrap_Returns_A_Handler_That_Sets_The_Claims_In_The_Request_Context(t *testing.T) {
	middleware := NewMiddleware(
		WithKey([]byte("secret")),
	)

	token, err := gojwt.NewWithClaims(gojwt.SigningMethodHS256, &gojwt.RegisteredClaims{Subject: "test"}).SignedString([]byte("secret"))
	assert.NoError(t, err)

	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(ContextKeyClaims).(*gojwt.RegisteredClaims)
		assert.Equal(t, "test", claims.Subject)
	}))

	handler.ServeHTTP(nil, &http.Request{
		Header: http.Header{
			"Authorization": []string{"Bearer " + token},
		},
	})
}

func Test_That_Wrap_Returns_A_Handler_That_Does_Not_Set_The_Claims_In_The_Request_Context_When_The_Token_Is_Invalid(t *testing.T) {
	middleware := NewMiddleware(
		WithKey([]byte("secret")),
	)

	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(ContextKeyClaims)
		assert.Nil(t, claims)
	}))

	handler.ServeHTTP(nil, &http.Request{
		Header: http.Header{
			"Authorization": []string{"Bearer test"},
		},
	})
}

func Test_That_Wrap_Returns_A_Handler_That_Does_Not_Set_The_Claims_In_The_Request_Context_When_The_Token_Is_Valid_But_The_Audience_Is_Wrong(t *testing.T) {
	middleware := NewMiddleware(
		WithKey([]byte("secret")),
		WithAudience("test"),
	)

	token, err := gojwt.NewWithClaims(gojwt.SigningMethodHS256, &gojwt.RegisteredClaims{Subject: "test"}).SignedString([]byte("secret"))
	assert.NoError(t, err)

	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(ContextKeyClaims)
		assert.Nil(t, claims)
	}))

	handler.ServeHTTP(nil, &http.Request{
		Header: http.Header{
			"Authorization": []string{"Bearer " + token},
		},
	})
}

func Test_That_Wrap_Returns_A_Handler_That_Does_Not_Set_The_Claims_In_The_Request_Context_When_The_Token_Is_Valid_But_The_Issuer_Is_Wrong(t *testing.T) {
	middleware := NewMiddleware(
		WithKey([]byte("secret")),
		WithIssuer("test"),
	)

	token, err := gojwt.NewWithClaims(gojwt.SigningMethodHS256, &gojwt.RegisteredClaims{Subject: "test"}).SignedString([]byte("secret"))
	assert.NoError(t, err)

	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(ContextKeyClaims)
		assert.Nil(t, claims)
	}))

	handler.ServeHTTP(nil, &http.Request{
		Header: http.Header{
			"Authorization": []string{"Bearer " + token},
		},
	})
}

func Test_That_Wrap_Returns_A_Handler_That_Sets_The_Claims_In_The_Request_Context_When_The_Token_Is_Valid_And_The_Audience_And_Issuer_Are_Correct(t *testing.T) {
	middleware := NewMiddleware(
		WithKey([]byte("secret")),
		WithAudience("test"),
		WithIssuer("test"),
	)

	token, err := gojwt.NewWithClaims(gojwt.SigningMethodHS256, &gojwt.RegisteredClaims{Subject: "test", Audience: []string{"test"}, Issuer: "test"}).SignedString([]byte("secret"))
	assert.NoError(t, err)

	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(ContextKeyClaims).(*gojwt.RegisteredClaims)
		assert.Equal(t, "test", claims.Subject)
		assert.Equal(t, 1, len(claims.Audience))
		assert.Equal(t, "test", claims.Audience[0])
		assert.Equal(t, "test", claims.Issuer)
	}))

	handler.ServeHTTP(nil, &http.Request{
		Header: http.Header{
			"Authorization": []string{"Bearer " + token},
		},
	})
}

func Test_That_Wrap_Returns_A_Handler_That_Sets_The_Claims_In_The_Request_Context_When_A_KeyFunc_Is_Provided(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	middleware := NewMiddleware(
		WithKey(func(*gojwt.Token) (interface{}, error) {
			return &privateKey.PublicKey, nil
		}),
	)

	token, err := gojwt.NewWithClaims(gojwt.SigningMethodRS256, &gojwt.RegisteredClaims{Subject: "test"}).SignedString(privateKey)
	assert.NoError(t, err)

	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(ContextKeyClaims).(*gojwt.RegisteredClaims)
		assert.Equal(t, "test", claims.Subject)
	}))

	handler.ServeHTTP(nil, &http.Request{
		Header: http.Header{
			"Authorization": []string{"Bearer " + token},
		},
	})
}

func Test_That_Wrap_Returns_A_Handler_That_Sets_The_Claims_In_The_Request_Context_When_A_Niladic_KeyFunc_Is_Provided(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	assert.NoError(t, err)

	middleware := NewMiddleware(
		WithKey(func() (interface{}, error) {
			return &privateKey.PublicKey, nil
		}),
	)

	token, err := gojwt.NewWithClaims(gojwt.SigningMethodRS256, &gojwt.RegisteredClaims{Subject: "test"}).SignedString(privateKey)
	assert.NoError(t, err)

	handler := middleware.Wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value(ContextKeyClaims).(*gojwt.RegisteredClaims)
		assert.Equal(t, "test", claims.Subject)
	}))

	handler.ServeHTTP(nil, &http.Request{
		Header: http.Header{
			"Authorization": []string{"Bearer " + token},
		},
	})
}
