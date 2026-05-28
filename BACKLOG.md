# Local RAG Backlog

## Priority 1: Citations & Source Attribution
- [x] **Data Structure Update**: Modify `internal/store/qdrant.go` to return a `SearchResult` struct containing `Content`, `SourceID`, and `PageNumber`.
- [x] **Metadata Resolution**: Implement `GetSource` in SQLite and resolve metadata during the `ask` command.
- [x] **Prompt Enhancement**: Update the prompt builder in `cmd/ask.go` to prefix every context chunk with its source (e.g., `[Source: book.pdf, Page: 12]`).
- [ ] **UI Update**: Add a "Sources" section at the end of the `ask` output to list the unique books/pages used for the answer (Final Polish).

## Priority 2: User Experience (Loading State)
- [ ] **CLI Spinner**: Implement a simple concurrent spinner in `cmd/ask.go` using a goroutine.
- [ ] **Feedback**: Ensure the spinner starts when the LLM begins generating and stops/clears when the answer is ready to print.

## Priority 3: Architecture (Refactors & Principles)
- [x] **Singleton Config**: Centralized configuration with environment variable support.
- [ ] **Cooperative Signaling**: Propagate `context.Context` through all layers (Ingestor, AI, Store).
    - [x] **Heartbeat**: Implemented in `ExtractChunk`.
    - [ ] **Fail-Fast**: Implement in `BatchEmbed` using an **Error Channel** (Decision made).
    - [ ] **Store Propagation**: Update `qdrant.go` to use `context.Context`.

## Priority 4: Ingestion Quality (Advanced Chunking)
- [ ] **Statistical Semantic Chunking**: Implement a sentence-based splitter that uses cosine similarity between sentence embeddings to identify logical breakpoints.
- [ ] **Context Overlap**: Add a configurable sliding window overlap to prevent context loss at chunk boundaries.

## Future Ideas
- [ ] Support for EPUB files.
- [ ] Interactive "Chat" mode (persistent conversation history).
- [ ] Topic discovery/auto-routing based on question embeddings.
