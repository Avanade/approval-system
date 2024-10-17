package timedjobs

type TimedJobs interface {
	ReprocessFailedCallbacks()
}
