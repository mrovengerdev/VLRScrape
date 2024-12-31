package paginator

import (
	"time"

	"golang.org/x/net/context"
	"golang.org/x/time/rate"
)

type Paginator struct {
	Limiter *rate.Limiter
	Context context.Context
	Cancel  context.CancelFunc
}

// TODO: Rewrite so cancel occurs upon 100 requests.
// Enforces paginator that limits requests to 10 per second and cancels out if it takes longer than 30 seconds.
func RestAPIPaginator() (paginator *Paginator) {
	limiter := rate.NewLimiter(1, 1) // 10 requests per second, with a burst of 1 request.

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	return &Paginator{
		Limiter: limiter,
		Context: ctx,
		Cancel:  cancel,
	}
}
