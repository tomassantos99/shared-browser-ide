package storage

import (
	"fmt"
	"time"
	"github.com/sirupsen/logrus"
)

func MemoryCleanup(storage *SessionStorage, cleanupRateSeconds int32, allowedInactiveSeconds int32) {
	ticker := time.NewTicker(time.Duration(cleanupRateSeconds) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		storage.mut.Lock()

		for id, session := range storage.sessions {
			fmt.Println(session.Id.String() + " " + session.LastActive.String())
			if time.Since(session.LastActive) > time.Duration(allowedInactiveSeconds) * time.Second {
				logrus.Info(fmt.Sprintf("Found an inactive session for more than %d seconds: %s", allowedInactiveSeconds, id))
				session.Close <- true
			}
		}

		storage.mut.Unlock()
	}
}
