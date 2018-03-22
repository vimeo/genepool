package workpool

// JobFunc is a function closure of work to be done
type JobFunc func(interface{}) (interface{}, error)

// Results contain the data returned by the job and its error status
type Results struct {
	D interface{} // the data returned by the job
	E error       // the status returned by the job
}

// Work includes a function closure to perform work and a channel to indicate when the work is complete
type work struct {
	f    JobFunc      // function closure of work to be done
	arg  interface{}  // input (if any) to the function, optionally allowing it to be a function and not necessarily a closure
	done chan Results // indicates when work is complete along with any data returned
}

// newWork returns a Work struct containing a function closure ad a Done channel
func newWork(f JobFunc, arg interface{}) work {

	// `done` channel notifies us when the work is complete
	done := make(chan Results)

	// Work contains our job closure and a channel to indicate when work is done
	return work{f, arg, done}
}
