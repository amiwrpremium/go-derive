package auth_test

import "time"

// timeT is an alias for time.Time so other test files can declare
// `timeNowDeterministic() timeT` without importing time directly.
type timeT = time.Time
