package contextworkpool

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mkrufky/genepool/workpool"
	. "github.com/smartystreets/goconvey/convey"
)

func getSleepJobFunc() workpool.JobFunc {
	return func(i interface{}) (interface{}, error) {
		dur, ok := i.(time.Duration)
		if !ok {
			return nil, fmt.Errorf("bad input")
		}

		time.Sleep(dur)
		return i, nil
	}
}

func TestContextWorkpool(t *testing.T) {
	Convey("Given a Workpool, a function, and an argument, `Run()` runs the job syncronously, waits for the result and returns job results and status when complete or context is cancelled, whichever occurs first", t, func() {
		ctx := context.Background()
		Convey("When the function completes without error, `Run()` returns the correct data and nil error", func() {

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			workpool := New(1, 1)
			val := 10 * time.Millisecond
			result, err := workpool.Run(ctx, getSleepJobFunc(), val)

			So(result, ShouldEqual, val)
			So(err, ShouldBeNil)
		})
		Convey("When the function completes with error (ie: bad type assertion), `Run()` returns the appropriate error", func() {

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			workpool := New(1, 1)
			val := "lollipop"
			result, err := workpool.Run(ctx, getSleepJobFunc(), val)

			So(result, ShouldBeNil)
			So(err, ShouldBeError)
		})
		Convey("When the context is not cancelled, `Run()` returns the correct data and nil error", func() {

			ctx, cancel := context.WithCancel(ctx)
			defer cancel()

			wp := New(1, 1)
			val := 10 * time.Millisecond

			var result workpool.Results
			done := make(chan interface{})

			go func(t *testing.T) {
				Convey("`start()` returns the correct data and nil error", t, func() {

					res, err := wp.Run(ctx, getSleepJobFunc(), val)

					So(res, ShouldEqual, val)
					So(err, ShouldBeNil)

					result = workpool.Results{D: res, E: err}
					done <- nil
				})
			}(t)

			// results should be empty until the job completes
			So(result.D, ShouldBeNil)
			So(result.E, ShouldBeNil)
			<-done

			// check again after job is complete
			So(result.D, ShouldEqual, val)
			So(result.E, ShouldBeNil)
		})
		Convey("When the context is cancelled before the function completes, `Run()` returns the appropriate error", func() {

			ctx, cancel := context.WithCancel(ctx)

			wp := New(1, 1)
			val := 10 * time.Millisecond

			var result workpool.Results
			done := make(chan interface{})

			go func(t *testing.T) {
				Convey("`start()` returns the appropriate error", t, func() {

					res, err := wp.Run(ctx, getSleepJobFunc(), val)

					So(res, ShouldBeNil)
					So(err, ShouldBeError)

					result = workpool.Results{D: res, E: err}
					done <- nil
				})
			}(t)

			cancel()

			So(result.D, ShouldBeNil)
			So(result.E, ShouldBeNil)
			<-done

			// check again after job is complete
			So(result.D, ShouldBeNil)
			So(result.E, ShouldBeError)
		})
	})
}
