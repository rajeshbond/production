// internal/contextkey/keys.go
package contextkey

// Define a custom string type for context keys to prevent key collisions
// with keys defined by other packages or the standard library.
type Key string

// KeyUser is the key used to store the authenticated user ID in the request context.
const KeyUser Key = "userID"

// You can add other context keys here if needed, e.g.:
// const KeyRequestID Key = "requestID"
