# Local RAG Backlog

## Priority 1: Citations & Source Attribution
- [ ] **Data Structure Update**: Modify `internal/store/qdrant.go` to return a `SourceChunk` struct containing `Content`, `FilePath`, and `PageNumber` instead of a raw string.
- [ ] **Prompt Enhancement**: Update the prompt builder in `cmd/ask.go` to prefix every context chunk with its source (e.g., `[Source: book.pdf, Page: 12]`).
- [ ] **UI Update**: Add a "Sources" section at the end of the `ask` output to list the unique books/pages used for the answer.

## Priority 2: User Experience (Loading State)
- [ ] **CLI Spinner**: Implement a simple concurrent spinner in `cmd/ask.go` using a goroutine.
- [ ] **Feedback**: Ensure the spinner starts when the LLM begins generating and stops/clears when the answer is ready to print.

## Priority 3: Architecture (Singleton Config)
- [x] **New Package**: Create `internal/config/config.go`.
- [x] **Singleton Pattern**: Use `sync.Once` to initialize a global configuration struct.
- [x] **Refactor**: Replace hardcoded literals in `ai/client.go`, `store/qdrant.go`, and `cmd/` with calls to `config.Get()`.

## Priority 4: Ingestion Quality (Advanced Chunking)
- [ ] **Statistical Semantic Chunking**: Implement a sentence-based splitter that uses cosine similarity between sentence embeddings to identify logical breakpoints.
- [ ] **Context Overlap**: Add a configurable sliding window overlap to prevent context loss at chunk boundaries.

## Future Ideas
- [ ] Support for EPUB files.
- [ ] Interactive "Chat" mode (persistent conversation history).
- [ ] Topic discovery/auto-routing based on question embeddings.
