package ezutil

type Job struct {
	logger      Logger
	setupFunc   func() error
	runFunc     func() error
	cleanupFunc func() error
}

func NewJob(logger Logger, runFunc func() error) *Job {
	if logger == nil {
		panic("logger cannot be nil")
	}
	if runFunc == nil {
		panic("runFunc cannot be nil")
	}
	return &Job{
		logger:  logger,
		runFunc: runFunc,
	}
}

func (j *Job) WithSetupFunc(fn func() error) *Job {
	j.setupFunc = fn
	return j
}

func (j *Job) WithCleanupFunc(fn func() error) *Job {
	j.cleanupFunc = fn
	return j
}

func (j *Job) Run() {
	if j.setupFunc != nil {
		j.logger.Info("setting up job...")
		if err := j.setupFunc(); err != nil {
			j.logger.Fatalf("error setting up job: %v", err)
		}
	}

	j.logger.Info("running job...")

	var jobErr error

	latency := MeasureLatency(func() { jobErr = j.runFunc() })

	if jobErr != nil {
		j.logger.Fatalf("error running job: %v", jobErr)
	}

	j.logger.Infof("success running job for %d ms", latency.Milliseconds())

	if j.cleanupFunc != nil {
		j.logger.Info("cleaning up job...")
		if err := j.cleanupFunc(); err != nil {
			j.logger.Fatalf("error cleaning up job: %v", err)
		}
	}
}
