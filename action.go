package main

import (
	"time"
)

type Action func() (bool, error)

var (
    curActionKey int64
)

type ActionQueue struct {
    actions map[int64]Action
}

func (a *ActionQueue) Update() error {
    finishedActions := make([]int64, 0)
    for key, action := range a.actions {
        if complete, err := action(); err != nil {
            return err
        } else if complete {
            finishedActions = append(finishedActions, key)
        }
    }
    for _, key := range finishedActions {
        delete(a.actions, key)
    }
    return nil
}

func (a *ActionQueue) Add(action Action) {
    a.actions[curActionKey] = action
    curActionKey++
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
        if time.Now().After(runTime) {
            return true, f()
        }
        return false, nil
    }
}
