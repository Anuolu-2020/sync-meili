package webhook

import "time"

type TaskDetails struct {
	UID               int       `json:"uid"`
	IndexUID          string    `json:"indexUid"`
	Status            string    `json:"status"`
	Type              string    `json:"type"`
	CanceledBy        *string   `json:"canceledBy"`
	ReceivedDocuments int       `json:"details.receivedDocuments"`
	IndexedDocuments  int       `json:"details.indexedDocuments"`
	Duration          string    `json:"duration"`
	EnqueuedAt        time.Time `json:"enqueuedAt"`
	StartedAt         time.Time `json:"startedAt"`
	FinishedAt        time.Time `json:"finishedAt"`
}
