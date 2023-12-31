package schedule

import (
	"context"
	"testing"
)

func TestFIFOSchedule(t *testing.T) {
	s := NewFIFOScheduler()
	defer s.Stop()

	next := 0
	jobCreator := func(i int) Job {
		return func(ctx context.Context) {
			if next != i {
				t.Fatalf("job#%d: got %d, want %d", i, next, i)
			}
			next = i + 1
		}
	}

	var jobs []Job
	for i := 0; i < 100; i++ {
		jobs = append(jobs, jobCreator(i))
	}

	for _, j := range jobs {
		s.Schedule(j)
	}

	s.WaitFinish(100)
	if s.Scheduled() != 100 {
		t.Errorf("scheduled = %d, want %d", s.Scheduled(), 100)
	}
}
