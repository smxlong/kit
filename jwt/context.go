package jwt

// ContextKey is a type for context keys.
type ContextKey string

const (
	// ContextKeyClaims is the context key for claims.
	ContextKeyClaims ContextKey = "jwtClaims"
)
