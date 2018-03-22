package workpool

import (
	"fmt"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func getSleepJobFunc() JobFunc {
	return func(i interface{}) (interface{}, error) {
		dur, ok := i.(time.Duration)
		if !ok {
			return nil, fmt.Errorf("bad input")
		}

		time.Sleep(dur)
		return i, nil
	}
}

func TestWorkpool(t *testing.T) {
	Convey("Given a Workpool, a function, and an argument, `start()` starts the job asyncronously and returns a channel that indicates when complete", t, func() {
		Convey("When the function completes without error, `start()` returns the correct data and nil error", func() {

			workpool := New(1, 1)
			var result Results
			done := make(chan interface{})
			val := 10 * time.Millisecond

			chResults := workpool.Start(getSleepJobFunc(), val)

			go func(t *testing.T) {
				Convey("`start()` returns the correct data and nil error", t, func() {
					result = <-chResults
					So(result.D, ShouldEqual, val)
					So(result.E, ShouldBeNil)
					done <- nil
				})
			}(t)

			<-done
		})
		Convey("When the function completes with error (ie: bad type assertion), `start()` returns the appropriate error", func() {

			workpool := New(1, 1)
			var result Results
			done := make(chan interface{})
			val := "lollipop"

			chResults := workpool.Start(getSleepJobFunc(), val)

			go func(t *testing.T) {
				Convey("`start()` returns the appropriate error", t, func() {
					result = <-chResults
					So(result.D, ShouldBeNil)
					So(result.E, ShouldBeError)
					done <- nil
				})
			}(t)

			<-done
		})
	})
	Convey("Given a Workpool, a function, and an argument, `Run()` runs the job syncronously, waits for the result and returns job results and status when complete", t, func() {
		Convey("When the function completes without error, `Run()` returns the correct data and nil error", func() {

			workpool := New(1, 1)
			val := 10 * time.Millisecond
			result, err := workpool.Run(getSleepJobFunc(), val)

			So(result, ShouldEqual, val)
			So(err, ShouldBeNil)
		})
		Convey("When the function completes with error (ie: bad type assertion), `Run()` returns the appropriate error", func() {

			workpool := New(1, 1)
			val := "lollipop"
			result, err := workpool.Run(getSleepJobFunc(), val)

			So(result, ShouldBeNil)
			So(err, ShouldBeError)
		})
	})
}
