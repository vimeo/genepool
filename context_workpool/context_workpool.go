package contextworkpool

import "github.com/mkrufky/workpool"
import "context"

type ContextWorkpool struct {
	*workpool.Workpool
}

func jobWithContext(ctx context.Context, f workpool.JobFunc) workpool.JobFunc {

	return func(i interface{}) (interface{}, error) {

		// first ensure that we still have context
		err := ctx.Err()
		if err != nil {
			return nil, err
		}

		// now call the function
		return f(i)
	}
}

// Run feeds work to the workpool and notifies the caller of error status when complete
// The supplied context is checked before each job
func (w *ContextWorkpool) Run(ctx context.Context, f workpool.JobFunc, arg interface{}) (interface{}, error) {

	// augment the supplied function to check the context before the execution
	jwc := jobWithContext(ctx, f)

	// start work and wait for it to complete, with early return on context error
	select {
	case <-ctx.Done():
		// return the context error status
		return nil, ctx.Err()
	case done := <-w.Start(jwc, arg):
		// return the job data and error status
		return done.D, done.E
	}
}

func New(workerCount, chanBuffer int) *ContextWorkpool {
	return &ContextWorkpool{workpool.New(workerCount, chanBuffer)}
}
