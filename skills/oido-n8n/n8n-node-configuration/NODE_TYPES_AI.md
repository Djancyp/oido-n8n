# AI & LangChain Nodes

LangChain nodes use the prefix `@n8n/n8n-nodes-langchain.*` (not `n8n-nodes-base.*`).

---

## Agents & Chains

| Display Name | Type | typeVersion | Notes |
|---|---|---|---|
| AI Agent | `@n8n/n8n-nodes-langchain.agent` | 1 | Main agent node. Needs model + optional tools/memory |
| Basic LLM Chain | `@n8n/n8n-nodes-langchain.chainLlm` | 1 | Single prompt → LLM → output |
| Summarization Chain | `@n8n/n8n-nodes-langchain.chainSummarization` | 2 | Summarize long documents |
| Retrieval QA Chain | `@n8n/n8n-nodes-langchain.chainRetrievalQa` | 1 | RAG: query + vector store |

## Language Models

| Display Name | Type | typeVersion | Credential |
|---|---|---|---|
| OpenAI Chat Model | `@n8n/n8n-nodes-langchain.lmChatOpenAi` | 1 | `openAiApi` |
| Anthropic Chat Model | `@n8n/n8n-nodes-langchain.lmChatAnthropic` | 1 | `anthropicApi` |
| Google Gemini Chat | `@n8n/n8n-nodes-langchain.lmChatGoogleGemini` | 1 | `googlePalmApi` |
| Mistral Cloud Chat | `@n8n/n8n-nodes-langchain.lmChatMistralCloud` | 1 | `mistralCloudApi` |
| Ollama Chat | `@n8n/n8n-nodes-langchain.lmChatOllama` | 1 | `ollamaApi` |
| Azure OpenAI Chat | `@n8n/n8n-nodes-langchain.lmChatAzureOpenAi` | 1 | `azureOpenAiApi` |
| AWS Bedrock Chat | `@n8n/n8n-nodes-langchain.lmChatAwsBedrock` | 1 | `aws` |
| Groq | `@n8n/n8n-nodes-langchain.lmChatGroq` | 1 | `groqApi` |
| Cohere | `@n8n/n8n-nodes-langchain.lmCohere` | 1 | `cohereApi` |

## Memory

| Display Name | Type | typeVersion | Notes |
|---|---|---|---|
| Window Buffer Memory | `@n8n/n8n-nodes-langchain.memoryBufferWindow` | 1 | Last N messages in memory |
| Postgres Chat Memory | `@n8n/n8n-nodes-langchain.memoryPostgresChat` | 1 | Persistent memory in Postgres |
| Redis Chat Memory | `@n8n/n8n-nodes-langchain.memoryRedisChat` | 1 | Persistent memory in Redis |
| MongoDB Chat Memory | `@n8n/n8n-nodes-langchain.memoryMongoDbChat` | 1 | Persistent memory in MongoDB |
| Motorhead | `@n8n/n8n-nodes-langchain.memoryMotorhead` | 1 | Motorhead memory server |
| Xata | `@n8n/n8n-nodes-langchain.memoryXata` | 1 | Xata memory store |

## Tools (usable by Agent)

| Display Name | Type | typeVersion | Notes |
|---|---|---|---|
| HTTP Request Tool | `@n8n/n8n-nodes-langchain.toolHttpRequest` | 1 | Agent calls an HTTP endpoint |
| Code Tool | `@n8n/n8n-nodes-langchain.toolCode` | 1 | Agent runs JS/Python code |
| Workflow Tool | `@n8n/n8n-nodes-langchain.toolWorkflow` | 1 | Agent calls another workflow |
| Calculator | `@n8n/n8n-nodes-langchain.toolCalculator` | 1 | Math expressions |
| SerpAPI (web search) | `@n8n/n8n-nodes-langchain.toolSerpApi` | 1 | Google search |
| Wikipedia | `@n8n/n8n-nodes-langchain.toolWikipedia` | 1 | Wikipedia lookups |
| Vector Store Tool | `@n8n/n8n-nodes-langchain.toolVectorStore` | 1 | Query a vector store |

## Vector Stores

| Display Name | Type | typeVersion | Notes |
|---|---|---|---|
| Pinecone | `@n8n/n8n-nodes-langchain.vectorStorePinecone` | 1 | Pinecone vector DB |
| Qdrant | `@n8n/n8n-nodes-langchain.vectorStoreQdrant` | 1 | Qdrant vector DB |
| Supabase | `@n8n/n8n-nodes-langchain.vectorStoreSupabase` | 1 | pgvector on Supabase |
| Postgres (pgvector) | `@n8n/n8n-nodes-langchain.vectorStorePostgres` | 1 | pgvector extension |
| In-Memory | `@n8n/n8n-nodes-langchain.vectorStoreInMemory` | 1 | Ephemeral, testing only |
| Redis | `@n8n/n8n-nodes-langchain.vectorStoreRedis` | 1 | Redis vector search |
| Zep | `@n8n/n8n-nodes-langchain.vectorStoreZep` | 1 | Zep memory server |

## Embeddings

| Display Name | Type | typeVersion | Credential |
|---|---|---|---|
| OpenAI Embeddings | `@n8n/n8n-nodes-langchain.embeddingsOpenAi` | 1 | `openAiApi` |
| Azure OpenAI Embeddings | `@n8n/n8n-nodes-langchain.embeddingsAzureOpenAi` | 1 | `azureOpenAiApi` |
| Google Vertex Embeddings | `@n8n/n8n-nodes-langchain.embeddingsGoogleVertex` | 1 | `googleApi` |
| Ollama Embeddings | `@n8n/n8n-nodes-langchain.embeddingsOllama` | 1 | `ollamaApi` |
| Cohere Embeddings | `@n8n/n8n-nodes-langchain.embeddingsCohere` | 1 | `cohereApi` |

## Document Loaders

| Display Name | Type | typeVersion | Notes |
|---|---|---|---|
| Default Data Loader | `@n8n/n8n-nodes-langchain.documentDefaultDataLoader` | 1 | Load from n8n items |
| Binary Input Loader | `@n8n/n8n-nodes-langchain.documentBinaryInputLoader` | 1 | Load from binary (PDF, DOCX) |
| GitHub Loader | `@n8n/n8n-nodes-langchain.documentGithubLoader` | 1 | Load from GitHub repo |

## Text Splitters

| Display Name | Type | typeVersion | Notes |
|---|---|---|---|
| Recursive Character Splitter | `@n8n/n8n-nodes-langchain.textSplitterRecursiveCharacterTextSplitter` | 1 | General purpose |
| Token Splitter | `@n8n/n8n-nodes-langchain.textSplitterTokenSplitter` | 1 | Split by token count |
| Character Splitter | `@n8n/n8n-nodes-langchain.textSplitterCharacterTextSplitter` | 1 | Split by character |

## Output Parsers

| Display Name | Type | typeVersion | Notes |
|---|---|---|---|
| Auto-fixing Output Parser | `@n8n/n8n-nodes-langchain.outputParserAutofixing` | 1 | Retry on parse failure |
| Structured Output Parser | `@n8n/n8n-nodes-langchain.outputParserStructured` | 1 | Force JSON schema output |
| Item List Output Parser | `@n8n/n8n-nodes-langchain.outputParserItemList` | 1 | Parse comma-separated list |

---

## AI Agent Workflow Pattern

```
Chat Trigger / Manual Trigger
  → AI Agent  ← Language Model (lmChatOpenAi)
               ← Memory (memoryBufferWindow)
               ← Tools: [toolHttpRequest, toolCode, toolWorkflow ...]
  → Output / Respond to Webhook
```

Node connections for AI sub-nodes use type `ai_languageModel`, `ai_memory`, `ai_tool` — not `main`:
```json
"connections": {
  "OpenAI Chat Model": {
    "ai_languageModel": [[{ "node": "AI Agent", "type": "ai_languageModel", "index": 0 }]]
  },
  "Window Buffer Memory": {
    "ai_memory": [[{ "node": "AI Agent", "type": "ai_memory", "index": 0 }]]
  },
  "HTTP Request Tool": {
    "ai_tool": [[{ "node": "AI Agent", "type": "ai_tool", "index": 0 }]]
  }
}
```
