package worker

import (
	"github.com/moira-alert/moira"
	"sync"
	"time"
)

func (worker *Checker) metricsChecker(metricEventsChannel <-chan *moira.MetricEvent) error {
	var handleWaitGroup sync.WaitGroup
	for {
		metricEvent, ok := <-metricEventsChannel
		if !ok {
			handleWaitGroup.Wait()
			worker.Logger.Info("Checking for new event stopped")
			return nil
		}
		handleWaitGroup.Add(1)
		go func(event *moira.MetricEvent) {
			defer handleWaitGroup.Done()
			if err := worker.handleMetricEvent(metricEvent); err != nil {
				worker.Logger.Errorf("Failed to handle metricEvent: %s", err.Error())
			}
		}(metricEvent)
	}
}

func (worker *Checker) handleMetricEvent(metricEvent *moira.MetricEvent) error {
	worker.lastData = time.Now().UTC().Unix()
	pattern := metricEvent.Pattern
	metric := metricEvent.Metric

	if err := worker.Database.AddPatternMetric(pattern, metric); err != nil {
		return err
	}
	triggerIds, err := worker.Database.GetPatternTriggerIDs(pattern)
	if err != nil {
		return err
	}
	if len(triggerIds) == 0 {
		if err := worker.Database.RemovePatternWithMetrics(pattern); err != nil {
			return err
		}
	}
	var performWaitGroup sync.WaitGroup
	worker.perform(triggerIds, worker.noCache, worker.Config.CheckInterval, &performWaitGroup)
	performWaitGroup.Wait()
	return nil
}
