package sdk

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"sync/atomic"
)

type decision struct {
	Request  CheckRequest
	Response CheckResponse
}

// CheckBatch - Performs a batch of authorization requests using Cedar Agent.
func (c Client) CheckBatch(
	ctx context.Context,
	reqs []CheckRequest,
	// TODO functional options
	numWorkers int,
) (map[CheckRequest]CheckResponse, error) {
	g, ctx := errgroup.WithContext(ctx)

	reqChan := make(chan CheckRequest)

	// STEP 1: Produce
	// We're feeding the "jobs" into our channel, and they'll be buffered and
	// picked up by our workers in the worker pool as soon as possible.
	g.Go(func() error {
		defer close(reqChan)
		for _, req := range reqs {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case reqChan <- req:
			}
		}
		return nil
	})

	decisions := make(chan *decision)

	// Step 2: Map
	// For each job, we'll perform the authorization request and write both the
	// request and the decision into a channel.
	workers := int32(numWorkers)
	for i := 0; i < numWorkers; i++ {
		g.Go(func() error {
			defer func() {
				// Last one out closes shop
				if atomic.AddInt32(&workers, -1) == 0 {
					close(decisions)
				}
			}()

			for req := range reqChan {
				if res, err := c.Check(ctx, req); err != nil {
					return fmt.Errorf(
						"unable to authorize principal %s for action %s on resource %s: %w",
						req.Principal, req.Action, req.Resource, err,
					)
				} else {
					select {
					case <-ctx.Done():
						return ctx.Err()
					case decisions <- &decision{
						Request:  req,
						Response: *res,
					}:
					}
				}
			}
			return nil
		})
	}

	// Step 3: Reduce
	// Transform decisions into the final output.
	ret := map[CheckRequest]CheckResponse{}
	g.Go(func() error {
		for d := range decisions {
			ret[d.Request] = d.Response
		}
		return nil
	})

	return ret, g.Wait()
}
