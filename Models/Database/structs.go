package Database

import (
	"sync"
	"time"
)

var (
	TotalInserts int
	StartTime    time.Time
	Mutex        sync.Mutex
)
