package sync

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/lib/pq"

	"github.com/Anuolu-2020/sync-meili/internal/config"
	"github.com/Anuolu-2020/sync-meili/internal/db"
	"github.com/Anuolu-2020/sync-meili/pkg/env"
)

func CreateTrigger(config config.SyncMeiliConfig) {
	_, err := db.DBManager.Conn.Exec(`
		CREATE OR REPLACE FUNCTION notify_change() RETURNS TRIGGER AS $$
		DECLARE
			data json;
		BEGIN
			IF (TG_OP = 'DELETE') THEN
				data = row_to_json(OLD);
			ELSE
				data = row_to_json(NEW);
			END IF;
			PERFORM pg_notify('table_changes', json_build_object('table', TG_TABLE_NAME, 'action', TG_OP, 'data', data)::text);
			RETURN NULL;
		END;
		$$ LANGUAGE plpgsql;
	`)
	if err != nil {
		log.Fatalf("Failed to create trigger function: %v", err)
	} else {
		log.Print("Trigger function created successfully")
	}

	for _, mapping := range config.Sync.Mappings {
		// Create the trigger
		_, err = db.DBManager.Conn.Exec(fmt.Sprintf(`
		CREATE TRIGGER %s_change_trigger
		AFTER INSERT OR UPDATE OR DELETE ON %s
		FOR EACH ROW EXECUTE FUNCTION notify_change();
      `, mapping.DatabaseTable, mapping.DatabaseTable))
		if err != nil {
			if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "42710" {
				fmt.Printf(
					"Trigger for table %s already exists, skipping.\n",
					mapping.DatabaseTable,
				)
				continue
			} else {
				log.Printf("Failed to create trigger for table %s: %v", mapping.DatabaseTable, err)
			}
		}
	}
}

func ListenAndSync(config *config.SyncMeiliConfig) {
	dbConnStr := env.GetEnv("DB_CONNECTION_STRING")
	// Create listener
	listener := pq.NewListener(dbConnStr, 10*time.Second, time.Minute, nil)
	defer listener.Close()

	err := listener.Listen("table_changes")
	if err != nil {
		log.Fatalf("Error while listening for table changes: %v", err)
	}

	fmt.Println("Listening for changes...")

	for {
		select {
		case <-time.After(90 * time.Second):
			go func() {
				listener.Ping() // Ping database every 90 seconds
			}()
		case notification := <-listener.Notify: // listen for changes in the database

			var payload map[string]interface{}
			// Unmarshal json from database
			err = json.Unmarshal([]byte(notification.Extra), &payload)
			if err != nil {
				log.Printf("Error unmarshaling pg notify JSON: %v", err)
				continue
			}

			// fmt.Printf("DB Payload: %v\n", payload)

			// Sync db row to meilisearch document
			SyncToMeilisearch(config, payload)

		}
	}
}
