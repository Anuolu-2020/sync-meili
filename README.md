# Syncmelli

Syncmelli is an experimental tool written in Go designed to synchronize data between MySQL and PostgreSQL databases with [Meilisearch](https://www.meilisearch.com/). It provides a simple and efficient way to keep your Meilisearch indexes up-to-date with changes in your relational databases.

## Features

- **Sync**: Sync data from MySQL and PostgreSQL to Meilisearch.
- **Real-time Updates**: Automatically update Meilisearch indexes when changes occur in your database.
- **Customizable Mapping**: Define how your database tables and columns map to Meilisearch documents.
- **Batch Processing**: Efficiently sync data in batches to reduce overhead.
- **Retry Mechanism**: Handle sync failures gracefully with configurable retries.

## Installation

To use Syncmelli, you need to have Go 1.16 or higher installed.

1. Clone the repository:
   ```bash
   git clone https://github.com/your-username/syncmelli.git
   cd syncmelli
   go build -o syncmelli
   ```

2. Configure Syncmelli by editing the config.yaml file (see Configuration below).

3. Run Syncmelli:
   ```bash
     ./syncmelli
   ```
## Configuration

Syncmelli uses a config.yaml file to define the synchronization settings. Below is an example configuration:

```yaml
    database:
  type: "postgres" # "mysql" or "postgres"
  server_id: 1 # for MySQL only

sync:
  enabled: true
  sync_request:
    max_retries: 3 # max retries on sync failure (e.g., connection issues/errors)
    retry_delay: "3s" # time to wait between retries
    request_timeout: "1s" # execution time for each sync request
  batch:
    batch_size: 4 # batch size before offloading
    flush_duration: "5s" # time to wait before offloading batch
  events:
    event_types:
      - "insert"
      - "update"
      - "delete"
  mappings:
    - database_table: "users" # defaults to public schema, for MySQL: "default.users" = db[default] table[users]
      meilisearch_index: "users"
      meilisearch_index_document_uid: "id"
      fields:
        - id: ""
        - name: ""
        - email: ""
    - database_table: "products"
      meilisearch_index: "products"
      meilisearch_index_document_uid: "id"
      fields:
        - "*"

```
## Configuration Options
  Database:
    type: The type of database (mysql or postgres).
    server_id: Required for MySQL only. Used for replication purposes.
 
  Sync:
    enabled: Enable or disable synchronization.
    sync_request:
      max_retries: Maximum number of retries on sync failure.
      retry_delay: Time to wait between retries (e.g., "3s").
      request_timeout: Maximum execution time for each sync request (e.g., "1s").
    batch:
      batch_size: Number of records to batch before offloading to Meilisearch.
      flush_duration: Time to wait before offloading a batch (e.g., "5s").
     events:
       event_types: List of database events to sync (insert, update, delete).
     mappings:
       database_table: The database table to sync (e.g., "users" or "default.users" for MySQL).
       meilisearch_index: The Meilisearch index name.
       meilisearch_index_document_uid: The primary key field in the Meilisearch index.
       fields: List of fields to sync. Use "*" to sync all fields.

## Usage
Once configured, Syncmelli will automatically sync the specified tables from MySQL or PostgreSQL to Meilisearch. You can start the sync process by running:
```bash
  ./syncmelli

```
Syncmelli will monitor changes in your database and update the corresponding Meilisearch indexes in real-time.

### Note: 
 Syncmelli is currently in an experimental phase. Use it at your own risk.
