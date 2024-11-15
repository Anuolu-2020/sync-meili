package sync

import (
"time"

"github.com/failsafe-go/failsafe-go/retrypolicy")

retryPolicy := retrypolicy.Builder[any]().WithDelay(time.Second).WithMaxRetries(3).Build()
