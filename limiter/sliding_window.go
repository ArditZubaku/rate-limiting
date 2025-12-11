package limiter

import (
	"container/list"
	"time"
)

type SlidingWindow struct {
	window int64
	limit  int
	logs   *list.List // deque // push_back, push_front -> O(1)
}

var _ RateLimiter = (*SlidingWindow)(nil)

func NewSlidingWindow(window int64, limit int) *SlidingWindow {
	return &SlidingWindow{
		window: window,
		limit:  limit,
		logs:   list.New(),
	}
}

func (s *SlidingWindow) Allow() bool {
	now := time.Now()
	edgeTime := time.Unix(now.Unix()-s.window, 0)

	// Remove outdated logs
	for s.logs.Len() > 0 {
		frontElement := s.logs.Front()
		if frontElement.Value.(time.Time).Before(edgeTime) {
			s.logs.Remove(frontElement)
		} else {
			break
		}
	}

	// Check if we can accept the request
	if s.logs.Len() < s.limit {
		s.logs.PushBack(now)
		return true // Accept request
	}

	return false // Reject request
}
