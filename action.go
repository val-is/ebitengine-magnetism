package main

import (
	"time"
)

type Action func() (bool, error)

type ActionQueue struct {
    actions []Action
}

func (a *ActionQueue) Update() error {
    updatedActions := make([]Action, 0)
    for _, action := range a.actions {
        if complete, err := action(); err != nil {
            return err
        } else if !complete {
            updatedActions = append(updatedActions, action)
        }
    }
    a.actions = updatedActions
    return nil
}

func (a *ActionQueue) Add(action Action) {
    a.actions = append(a.actions, action)
}

func NewContinuousTimedAction(f func(percentComplete float64, duration time.Duration) (doneEarly bool, err error), duration time.Duration) Action {
    startTime := time.Now()
    endTime := startTime.Add(duration)
    return func() (bool, error) {
        curTime := time.Now()
        if curTime.After(endTime) {
            return true, nil
        }
        percentComplete := curTime.Sub(startTime).Seconds() / duration.Seconds()
        if doneEarly, err := f(percentComplete, duration); err != nil {
            return false, err
        } else if doneEarly {
            return true, nil
        }
        return false, nil
    }
}

func NewTimerAction(f func() error, runTime time.Time) Action {
    return func() (bool, error) {
        if time.Now().Before(runTime) {
            return false, nil
        }
        return true, f()
    }
}
