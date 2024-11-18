package sync

import (
	"fmt"
	"log"
	"time"

	"github.com/failsafe-go/failsafe-go"
	"github.com/failsafe-go/failsafe-go/retrypolicy"
	"github.com/failsafe-go/failsafe-go/timeout"
	"github.com/meilisearch/meilisearch-go"

	"github.com/Anuolu-2020/sync-meili/internal/config"
	"github.com/Anuolu-2020/sync-meili/pkg/env"
)

type MeiliSearchRequest struct {
	DocumentId       string
	Action           string
	MeiliSearchIndex string
	Data             map[string]interface{}
}

var (
	INSERT_DOCUMENT = "INSERT"
	UPDATE_DOCUMENT = "UPDATE"
	DELETE_DOCUMENT = "DELETE"
)

// Configure meilisearch client and return the client
func configureMeilisearchClient() meilisearch.ServiceManager {
	client := meilisearch.New(
		env.GetEnv("MEILISEARCH_CONN_STRING"),
		meilisearch.WithAPIKey(env.GetEnv("MEILISEARCH_API_KEY")),
	)

	return client
}

// Get meilisearch index from client
func GetMeilisearchIndex(
	client meilisearch.ServiceManager,
	indexUid string,
) meilisearch.IndexManager {
	index := client.Index(indexUid)

	return index
}

func SyncInBatchRequest(
	input <-chan MeiliSearchRequest,
	config *config.SyncMeiliConfig,
) {
	maxRetries := config.Sync.SyncRequest.MaxRetries
	if maxRetries == 0 {
		maxRetries = 3 // Defaults to 3 retries
	}

	retryDelay := config.Sync.SyncRequest.RetryDelay
	if int64(retryDelay) == 0 {
		retryDelay = 2 * time.Second // defaults to 2 seconds
	}

	requestTimeout := config.Sync.SyncRequest.RequestTimeout
	if int64(requestTimeout) == 0 {
		requestTimeout = 1 * time.Second // defaults to 1 second
	}

	// Configure timeout policy for request
	timeoutPolicy := timeout.With[any](requestTimeout)

	// Configure retry policy for request
	retryPolicy := retrypolicy.Builder[any]().WithDelay(retryDelay).
		WithMaxRetries(maxRetries).
		AbortOnErrors(meilisearch.ErrRequestBodyWithoutContentType, meilisearch.ErrInvalidRequestMethod).
		OnRetryScheduled(func(e failsafe.ExecutionScheduledEvent[any]) {
			fmt.Println("Ping retry", e.Attempts(), "after delay of", e.Delay)
		}).Build()

	batchSize := config.Sync.Batch.BatchSize

	if batchSize == 0 {
		batchSize = 1 // Defaults to 1 meaning realtime sync
	}

	flush := config.Sync.Batch.FlushDuration

	if int64(flush) == 0 {
		flush = 5000 * time.Millisecond // Defaults to 5 Seconds
	}

	timer := time.NewTimer(flush)
	defer timer.Stop()

	// Initialize meilisearch client
	client := configureMeilisearchClient()

	insertDocumentsBatch := make(map[string][]map[string]interface{})
	updateDocumentsBatch := make(map[string][]map[string]interface{})
	deleteDocumentsBatch := make(map[string][]string)

	for {
		select {
		case request := <-input:
			switch request.Action {
			case INSERT_DOCUMENT:

				// batch document to index key
				insertDocumentsBatch[request.MeiliSearchIndex] = append(
					insertDocumentsBatch[request.MeiliSearchIndex],
					request.Data,
				)

				if len(insertDocumentsBatch[request.MeiliSearchIndex]) >= batchSize {
					index := GetMeilisearchIndex(client, request.MeiliSearchIndex)

					log.Printf(
						"Batch complete offloading %d",
						len(insertDocumentsBatch[request.MeiliSearchIndex]),
					)

					go AddDocumentToMeilisearch(
						index,
						insertDocumentsBatch[request.MeiliSearchIndex], retryPolicy, timeoutPolicy)

					// Reset the batch
					insertDocumentsBatch[request.MeiliSearchIndex] = nil
				}

			case UPDATE_DOCUMENT:
				updateDocumentsBatch[request.MeiliSearchIndex] = append(
					updateDocumentsBatch[request.MeiliSearchIndex],
					request.Data,
				)

				if len(updateDocumentsBatch[request.MeiliSearchIndex]) >= batchSize {
					index := GetMeilisearchIndex(client, request.MeiliSearchIndex)

					go UpdateDocumentsInMeilisearch(
						index,
						updateDocumentsBatch[request.MeiliSearchIndex], retryPolicy, timeoutPolicy)

					// Reset the batch
					updateDocumentsBatch[request.MeiliSearchIndex] = nil
				}

			case DELETE_DOCUMENT:

				deleteDocumentsBatch[request.MeiliSearchIndex] = append(
					deleteDocumentsBatch[request.MeiliSearchIndex],
					request.DocumentId,
				)

				if len(deleteDocumentsBatch[request.MeiliSearchIndex]) >= batchSize {
					index := GetMeilisearchIndex(client, request.MeiliSearchIndex)

					log.Printf(
						"Batch complete offloading %d",
						len(deleteDocumentsBatch[request.MeiliSearchIndex]),
					)

					go DeleteDocumentsFromMeilisearch(
						index,
						deleteDocumentsBatch[request.MeiliSearchIndex], retryPolicy, timeoutPolicy)

					// Reset the batch
					deleteDocumentsBatch[request.MeiliSearchIndex] = nil
				}

			}
		case <-timer.C:
			// log.Println("Flushing batches on timer expiry")

			// Offload batch
			for meilisearchIndex, documents := range insertDocumentsBatch {
				if len(documents) > 0 {
					index := GetMeilisearchIndex(client, meilisearchIndex)
					log.Printf("offloading %d documents", len(documents))
					go AddDocumentToMeilisearch(index, documents, retryPolicy, timeoutPolicy)

					// Reset the batch
					insertDocumentsBatch[meilisearchIndex] = nil
				}
			}

			for meilisearchIndex, documents := range updateDocumentsBatch {
				if len(documents) > 0 {
					index := GetMeilisearchIndex(client, meilisearchIndex)
					go UpdateDocumentsInMeilisearch(index, documents, retryPolicy, timeoutPolicy)

					// Reset the batch
					updateDocumentsBatch[meilisearchIndex] = nil
				}
			}

			for meilisearchIndex, documents := range deleteDocumentsBatch {
				if len(documents) > 0 {
					index := GetMeilisearchIndex(client, meilisearchIndex)
					go DeleteDocumentsFromMeilisearch(index, documents, retryPolicy, timeoutPolicy)

					// Reset the batch
					deleteDocumentsBatch[meilisearchIndex] = nil
				}
			}

			timer.Reset(flush)

		}
	}
}

func AddDocumentToMeilisearch(
	index meilisearch.IndexManager,
	documents []map[string]interface{},
	retrypolicy retrypolicy.RetryPolicy[any],
	timeoutPolicy timeout.Timeout[any],
) {
	err := failsafe.Run(func() error {
		task, err := index.AddDocuments(documents)
		if err != nil {
			log.Printf("Error inserting document to meilisearch, retrying...: %v\n", err)
			return err
		}
		log.Printf("Add document taskuid: %d\n", task.TaskUID)

		return nil
	}, retrypolicy, timeoutPolicy)
	if err != nil {
		log.Printf("Failed to add documents after retries: %v\n", err)
	}
}

func UpdateDocumentsInMeilisearch(
	index meilisearch.IndexManager,
	documents []map[string]interface{},
	retrypolicy retrypolicy.RetryPolicy[any],
	timeoutPolicy timeout.Timeout[any],
) {
	err := failsafe.Run(func() error {
		task, err := index.UpdateDocuments(documents)
		if err != nil {
			fmt.Printf("Error updating document in meilisearch, retrying...: %v", err)
			return err
		}
		log.Printf("Update document taskuid: %d\n", task.TaskUID)
		return nil
	}, retrypolicy, timeoutPolicy)
	if err != nil {
		log.Printf("Failed to update documents after retries: %v\n", err)
	}
}

func DeleteDocumentsFromMeilisearch(
	index meilisearch.IndexManager,
	documentId []string,
	retrypolicy retrypolicy.RetryPolicy[any],
	timeoutPolicy timeout.Timeout[any],
) {
	err := failsafe.Run(func() error {
		task, err := index.DeleteDocuments(documentId)
		if err != nil {
			log.Printf("Error deleting document in meilisearch, retrying...: %v", err)
			return err
		}
		log.Printf("Delete document taskuid: %d\n", task.TaskUID)
		return nil
	}, retrypolicy, timeoutPolicy)
	if err != nil {
		log.Printf("Failed to delete documents after retries: %v\n", err)
	}
}

func FilterDocumentAndSync(
	config *config.SyncMeiliConfig,
	payload map[string]interface{},
	batchChannel chan<- MeiliSearchRequest,
) {
	document := make(map[string]interface{})

	// Sync mappings for each table
	mappings := config.Sync.Mappings

	for _, mapping := range mappings {
		if mapping.DatabaseTable == payload["table"] {

			for _, field := range mapping.Fields {
				if field == "*" || field == "" {
					document = payload["data"].(map[string]interface{})
					break
				} else {
					document[field] = payload["data"].(map[string]interface{})[field]
				}
			}

			// Get document id value
			documentId := fmt.Sprintf(
				"%v",
				document[mapping.MeilisearchIndexDocumentUid])

			// Pass request to batch channel
			batchChannel <- MeiliSearchRequest{
				DocumentId:       documentId,
				Action:           payload["action"].(string),
				MeiliSearchIndex: mapping.MeilisearchIndex,
				Data:             document,
			}

		}
	}
}
