package workpool

// Workpool is an object that will allocate a channel for scheduling async work
// and start up a predefined number of workers to consume work from this channel
type Workpool struct {
	c chan work
}

// worker processes incoming work and notifies the caller of error status when complete
func (w *Workpool) worker() {
	// range over the worker channel to call the function closure and return its status over the done channel
	for j := range w.c {
		d, err := j.f(j.arg)
		j.done <- Results{d, err}
	}
}

// New returns a new worker with its own channel with a buffer of size `buffer`
func New(workerCount, chanBuffer int) *Workpool {

	// c channel accepts incoming work
	c := make(chan work, chanBuffer)

	// w is the worker object to be returned
	w := &Workpool{c}

	// start up the actual workers, all consuming work from this one channel
	for wc := 0; wc < workerCount; wc++ {
		go w.worker()
	}

	return w
}

// Start creates new work, feeds it to the workpool and returns a channel to notify the caller when complete
func (w *Workpool) Start(f JobFunc, arg interface{}) chan Results {

	// create a new work structure along with the complete notification channel
	work := newWork(f, arg)

	// feed this work to our workers
	w.c <- work

	// return a channel to notify the caller when this job has completed
	return work.done
}

// Run feeds work to the workpool and notifies the caller of error status when complete
func (w *Workpool) Run(f JobFunc, arg interface{}) (interface{}, error) {

	// start work and wait for it to complete
	ret := <-w.Start(f, arg)

	// return the job data and error status
	return ret.D, ret.E
}
