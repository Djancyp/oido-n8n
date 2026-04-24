# Data & Storage Nodes

## Relational Databases

| Display Name | Type | typeVersion | Credential |
|---|---|---|---|
| Postgres | `n8n-nodes-base.postgres` | 2 | `postgres` |
| MySQL | `n8n-nodes-base.mySql` | 2 | `mySql` |
| MariaDB | `n8n-nodes-base.mariaDb` | 2 | `mariaDb` |
| SQLite | `n8n-nodes-base.sqlite` | 1 | None (local file) |
| Microsoft SQL | `n8n-nodes-base.microsoftSql` | 2 | `microsoftSql` |
| CockroachDB | `n8n-nodes-base.cockroachDb` | 1 | `cockroachDb` |
| QuestDB | `n8n-nodes-base.questDb` | 1 | `questDb` |

## NoSQL / Document

| Display Name | Type | typeVersion | Credential |
|---|---|---|---|
| MongoDB | `n8n-nodes-base.mongoDb` | 1 | `mongoDb` |
| CouchDB | `n8n-nodes-base.couchDb` | 1 | `couchDb` |
| Elasticsearch | `n8n-nodes-base.elasticsearch` | 1 | `elasticsearchApi` |
| Cassandra | `n8n-nodes-base.cassandra` | 1 | `cassandra` |

## Key-Value / Cache

| Display Name | Type | typeVersion | Credential |
|---|---|---|---|
| Redis | `n8n-nodes-base.redis` | 1 | `redis` |

## Cloud Databases / BaaS

| Display Name | Type | typeVersion | Credential |
|---|---|---|---|
| Supabase | `n8n-nodes-base.supabase` | 1 | `supabaseApi` |
| Baserow | `n8n-nodes-base.baserow` | 1 | `baserowApi` |
| NocoDB | `n8n-nodes-base.nocoDb` | 3 | `nocoDb` |
| Airtable | `n8n-nodes-base.airtable` | 2 | `airtableTokenApi` |

## Spreadsheets

| Display Name | Type | typeVersion | Credential |
|---|---|---|---|
| Google Sheets | `n8n-nodes-base.googleSheets` | 4 | `googleSheetsOAuth2Api` |
| Microsoft Excel | `n8n-nodes-base.microsoftExcel` | 2 | `microsoftExcelOAuth2Api` |
| Spreadsheet File | `n8n-nodes-base.spreadsheetFile` | 2 | None (binary) |

## Cloud File Storage

| Display Name | Type | typeVersion | Credential |
|---|---|---|---|
| Google Drive | `n8n-nodes-base.googleDrive` | 3 | `googleDriveOAuth2Api` |
| Dropbox | `n8n-nodes-base.dropbox` | 1 | `dropboxApi` |
| Box | `n8n-nodes-base.box` | 1 | `boxOAuth2Api` |
| OneDrive | `n8n-nodes-base.microsoftOneDrive` | 1 | `microsoftOneDriveOAuth2Api` |
| AWS S3 | `n8n-nodes-base.s3` | 1 | `s3` |
| Google Cloud Storage | `n8n-nodes-base.googleCloudStorage` | 1 | `googleCloudStorageOAuth2Api` |

## Local Files & Binary

| Display Name | Type | typeVersion | Notes |
|---|---|---|---|
| Read/Write File | `n8n-nodes-base.readWriteFile` | 1 | Local filesystem |
| Extract From File | `n8n-nodes-base.extractFromFile` | 1 | Binary → JSON (PDF, CSV, HTML, image) |
| Convert To File | `n8n-nodes-base.convertToFile` | 1 | JSON → binary |
| Spreadsheet File | `n8n-nodes-base.spreadsheetFile` | 2 | Read/write XLSX, CSV |
| Compression | `n8n-nodes-base.compression` | 1 | Zip / gzip / tar |
| FTP | `n8n-nodes-base.ftp` | 1 | FTP / SFTP |

## Knowledge & Docs

| Display Name | Type | typeVersion | Credential |
|---|---|---|---|
| Notion | `n8n-nodes-base.notion` | 2 | `notionApi` |
| Confluence | `n8n-nodes-base.confluence` | 1 | `confluenceApi` |
| Google Docs | `n8n-nodes-base.googleDocs` | 2 | `googleDocsOAuth2Api` |

## Message Queues

| Display Name | Type | typeVersion | Credential |
|---|---|---|---|
| AWS SQS | `n8n-nodes-base.awsSqs` | 1 | `aws` |
| AWS SNS | `n8n-nodes-base.awsSns` | 1 | `aws` |
| RabbitMQ | `n8n-nodes-base.rabbitmq` | 1 | `rabbitmq` |
| MQTT | `n8n-nodes-base.mqtt` | 1 | `mqtt` |
| Kafka | `n8n-nodes-base.kafka` | 1 | `kafka` |
