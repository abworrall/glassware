package sensor

import(
	"fmt"
	"time"

	"github.com/abworrall/glassware/pkg/event"
)

var	NumValsToAverage = 4

// A DropRestoreSensor receives a stream of ints as readings, and
// identifies when the value has "dropped" - and then when it becomes
// "restored".
//
// In practice, adjusting weights on top of an FSR and reading the
// voltage, we see a lot of drift in the value. E.g. a running average
// of 660, then a drop to 430, then a restore to 580. The
// DropRestoreSensor maintains a recent smoothed average over the past
// few readings, so it auto-adjusts as the FSR drifts.
//
// A Drop or Restore correspond to a shift in ~6.6% of the value -
// downwards for a Drop, upwards for a Restore. Restores and Drops can
// only happen in sequence, e.g. a Drop then a Restore then a Drop, etc.
type DropRestoreSensor struct {
	Name              string // Must be unique across all controllers, and stable over reboots
	PrevVals        []int	
	IsInDropStatus    bool
}

func NewDropRestoreSensor(name string) *DropRestoreSensor {
	return &DropRestoreSensor{Name: name, PrevVals: []int{}}
}

func (drs *DropRestoreSensor)String() string {
	return fmt.Sprintf("DropRestore(%s: prev=[%v], isInDrop=%v)", drs.Name, drs.PrevVals, drs.IsInDropStatus)
}

func (drs *DropRestoreSensor)GetName() string { return drs.Name }

func (drs *DropRestoreSensor)IsCalibrated() bool { return len(drs.PrevVals) >= NumValsToAverage }

func (drs *DropRestoreSensor)RecentAverage() int {
	tot := 0
	for _, val := range drs.PrevVals {
		tot += val
	}
	return tot / len(drs.PrevVals)
}

func (drs *DropRestoreSensor)ProcessNewReading(val int, evOut chan<- event.Event) {
	if !drs.IsCalibrated() {
		drs.PrevVals = append(drs.PrevVals, val)
		return
	}

	// Now compare the value to our recent rolling average
	recentAvg := drs.RecentAverage()
	threshold := recentAvg / 15                // a ~6.6% shift is significant
	delta := val - recentAvg                   // -ve delta means a drop

	//log.Printf("                                       "+
	//	"%s: val=%d, delta=% 3d (recent=%d, thresh=%d, isDropped=%v)\n",
	//	drs.GetName(), val, delta, drs.RecentAverage(), threshold, drs.IsInDropStatus)

	if !drs.IsInDropStatus && delta < (-1 * threshold) {
		drs.IsInDropStatus = true
		evOut <- event.Event{Kind:event.EvDrop, SourceName: drs.GetName(), Time: time.Now()}

	} else if drs.IsInDropStatus && delta > threshold {
		drs.IsInDropStatus = false
		evOut <- event.Event{Kind:event.EvRestore, SourceName: drs.GetName(), Time: time.Now()}

	} else {
		drs.PrevVals = append(drs.PrevVals[1:], val) // Update the rolling average, adjust to the value drift
	}
}
