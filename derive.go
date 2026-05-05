package goderive

// Version is the SDK semantic version. It is reported in the User-Agent header
// of REST requests and in the WebSocket connect handshake. The literal is
// maintained by release-please via the x-release-please-version annotation.
const Version = "0.2.0" // x-release-please-version

// UserAgent returns the default User-Agent string used by transports.
func UserAgent() string {
	return "go-derive/" + Version
}
