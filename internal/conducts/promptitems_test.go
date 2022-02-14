package conducts

import (
	"testing"
	"time"
)

func TestClosetRefreshTime(t *testing.T) {
	var priorTime = 5 * time.Minute
	var interval = 10 * time.Minute
	var closest = newClosetRefreshTime(priorTime, interval)

	closest.update(30 * time.Minute)
	if closest.Time != 30*time.Minute-priorTime {
		t.Errorf("unexpected value: %d", closest.Time)
	}
	if closest.next() != interval {
		t.Errorf("unexpected value: %d", closest.next())
	}

	closest.update(35 * time.Minute)
	if closest.Time != 30*time.Minute-priorTime {
		t.Errorf("unexpected value: %d", closest.Time)
	}

	closest.update(25 * time.Minute)
	if closest.Time != 25*time.Minute-priorTime {
		t.Errorf("unexpected value: %d", closest.Time)
	}

	closest.update(interval + priorTime + 1*time.Minute)
	if closest.next() != interval {
		t.Errorf("unexpected value: %d", closest.next())
	}

	closest.update(interval + priorTime - 1*time.Minute)
	if closest.next() != interval-1*time.Minute {
		t.Errorf("unexpected value: %d", closest.next())
	}
}
