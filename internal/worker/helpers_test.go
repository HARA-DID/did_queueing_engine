package worker_test

import "errors"

// errSentinel is a generic sentinel used in handler tests.
var errSentinel = errors.New("sentinel blockchain error")
