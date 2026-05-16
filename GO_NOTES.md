# Go Learning Notes: Building a Local RAG CLI

These notes track key Go concepts discovered while building the `local-rag` project.

## 1. Package-Level Visibility
- **The Concept**: All files in the same folder with the same `package name` (e.g., `package cmd`) share the same scope.
- **Visibility**: 
    - Variables starting with **lowercase** (e.g., `var rootCmd`) are **private** to the package.
    - Variables starting with **uppercase** (e.g., `var Execute`) are **exported** (public) and visible to other packages.
- **Why it matters**: You don't need to import your own files. `ingest.go` can use `rootCmd` defined in `root.go` automatically.

## 2. The `init()` Function
- **Reserved Function**: `init()` is a special function that runs automatically before `main()`.
- **CLI Pattern**: Used in Cobra to "register" subcommands. Each command file handles its own registration, keeping the project modular.
- **Execution Order**: The compiler finds all `init()` functions in a package and runs them at startup.

## 3. Pointers & Memory Addresses (`&` and `*`)
- **The `&` Operator**: Used to get the **memory address** of a variable.
- **Why use them for Flags?**: When we pass `&path` to `StringVarP`, we are giving the Cobra library the "address" of our variable. This allows Cobra to reach into our memory and update the value when the user types a flag, rather than just changing a local copy.

## 4. Structs as Commands
- **Composition**: In Go, we use `structs` (like `cobra.Command`) to define objects and their behavior.
- **Function Literals (Anons)**: The `Run` field in a command is often assigned an anonymous function: `func(cmd *cobra.Command, args []string) { ... }`.

## 5. Idiomatic CLI Structure
- **Centralization vs. Modularity**: Go favors modularity. Instead of one giant list of commands in `root.go`, each command registers itself. This prevents "bottleneck" files as the project scales.

## 6. Short Variable Declaration (`:=`)
- **Combined Action**: Declares a new variable and assigns it a value in one step.
- **Inference**: Go automatically determines the type based on the value on the right side.
- **Scope**: Only works inside functions. Top-level variables must use the `var` keyword.
- **Limitation**: You cannot use `:=` to assign to a variable that has already been declared in the same scope. Use `=` for subsequent assignments.

## 7. Multiple Return Values
- **Pattern**: Functions in Go frequently return more than one value, usually `(result, error)`.
- **Handling**: You must handle or explicitly ignore all returned values. Using `_` allows you to ignore a value you don't need.

## 8. Interfaces (The "Contract")
- **Definition**: A set of method signatures. It defines **behavior**, not data.
- **Implicit Satisfaction**: You don't "implement" an interface explicitly. If a struct has the methods defined in an interface, it satisfies it automatically.
- **The `io.Reader`**: One of the most common interfaces. It represents anything that can be read into a byte slice.

## 9. The `defer` Keyword
- **Purpose**: Schedules a function call to run immediately before the surrounding function returns.
- **Common Use**: Cleanup tasks like closing files (`f.Close()`) or unlocking mutexes.
- **LIFO Order**: If you use multiple `defers`, they run in "Last-In, First-Out" order.

## 10. Function Signatures & Returns
- **Strictness**: A function must return exactly the number and type of values defined in its signature.
- **Cobra Run vs RunE**: The standard `Run` function in Cobra returns nothing. If you need to return an error to the user, use `RunE`, which expects an `error` return value.

## 11. Built-in Functions
- **Generics**: Functions like `min()` and `max()` (added in Go 1.21) are built-in and work across different numeric types.
- **Safety**: Use `min(desired, len(thing))` to safely slice strings or slices without going "out of bounds."

## 12. Methods & Receivers
- **No Classes**: Go uses methods on structs instead of classes.
- **Receiver Syntax**: `func (p *Processor) DoSomething()` attaches a function to a struct. The `p` is the "receiver" (similar to `self` or `this`).
- **Pointer Receivers**: Using `*Processor` allows the method to modify the struct's data and is more memory-efficient for large structs.

## 13. High-Performance Strings (`strings.Builder`)
- **Immutability**: Go strings are immutable. Adding strings with `+` in a loop creates many temporary objects in memory.
- **Builder**: Use `strings.Builder` to efficiently accumulate large amounts of text. Use `.WriteString()` and then `.String()` at the end.

## 14. Runes vs. Bytes
- **Strings as Bytes**: In Go, `string` is a slice of bytes. `len(str)` returns the number of bytes, not characters.
- **Runes**: A `rune` (alias for `int32`) represents a single Unicode character.
- **Safety**: When slicing strings that might contain non-ASCII characters, convert the string to `[]rune` first: `runes := []rune(str)`.

## 15. JSON Struct Tags
- **The Concept**: Backticks `` `json:"key_name"` `` tell Go how to rename struct fields when converting to/from JSON.
- **Why**: Go fields often start with Uppercase (for visibility), but JSON APIs often expect lowercase.

## 16. The "New..." Constructor Pattern
- **Convention**: Since Go lacks explicit constructors, we use functions like `func NewClient(...) *Client`.
- **Pointers**: Usually returns a pointer (`*`) to the struct so it can be shared and modified across the app.

## 17. Creating Custom Errors (`fmt.Errorf`)
- **Usage**: Use `fmt.Errorf("message: %v", value)` to create a new error on the fly.
- **Wrapping**: Use the `%w` verb to wrap an existing error, preserving the original error's context: `fmt.Errorf("context: %w", err)`.

## 18. Context (`context.Context`)
- **Purpose**: Carries deadlines, cancellation signals, and other request-scoped values across API boundaries and goroutines.
- **Convention**: Always pass it as the first argument in functions that perform I/O (Database calls, HTTP requests).
- **Background**: `context.Background()` is the default "empty" context used when you don't have a specific deadline yet.

## 19. Type Casting (Conversions)
- **Strictness**: Go does NOT implicitly convert types.
- **Syntax**: `destinationValue := typeName(sourceValue)`.

## 20. For-Range Loop Behavior (The "Copy" Gotcha)
- **Value Copy**: By default, `for i, v := range slice` creates a copy.
- **Modifying Slices**: Use the index: `slice[i].Field = value`.

## 21. Error Analysis & String Inspection
- **`err.Error()`**: Converts an error object into its string representation.
- **`strings.Contains(str, substr)`**: A safe way to check for specific keywords in an error message.

## 22. Goroutines (Concurrency)
- **The `go` Keyword**: Prefixing a function call with `go` starts a lightweight thread (goroutine) that runs independently.
- **`sync.WaitGroup`**: Used to wait for a collection of goroutines to finish. Use `.Add(n)`, `.Done()`, and `.Wait()`.

## 23. Channels (The Pipeline)
- **The Concept**: A typed conduit through which you can send and receive values.
- **Unbuffered**: `make(chan int)`. Sending and receiving block until both sides are ready (a "handshake").
- **Buffered**: `make(chan int, capacity)`. Sends only block when the buffer is full. Receives only block when the buffer is empty.
- **Closing**: Use `close(ch)` to signal completion. You can still receive from a closed channel until the buffer is empty, but you cannot send to it.

## 24. Package-Level Redeclaration Error
- **The Conflict**: Since all files in a folder share the same `package` scope, you cannot declare the same variable at the top level in two different files (e.g., `var topic string` in both `ingest.go` and `ask.go`).
- **The Error**: `topic redeclared in this block`.
- **Solution**: Declare shared variables once in a common file like `root.go`, or keep them local to a function if they don't need package-wide visibility.

## 25. Dependency Management
- **`go get <package>`**: Downloads a specific package and adds it to your `go.mod` file.
- **`go mod tidy`**: Cleans up your `go.mod` and `go.sum` files. It removes unused dependencies and adds missing ones based on your imports.
- **`go mod download`**: Downloads all dependencies listed in your `go.mod` file to your local cache. Useful for CI/CD or fresh setups.

## 26. Zero Values & Struct Initialization
- **The Concept**: In Go, variables are never "uninitialized" or `undefined`. If you don't provide a value, Go assigns a "Zero Value":
    - `int/float`: `0`
    - `string`: `""` (empty string)
    - `bool`: `false`
    - `pointers/slices/maps/channels`: `nil`
- **Structs as "Buckets"**: You can create a struct and only fill a few fields. The rest will hold their zero values. This is common when:
    1. **Preparing Data**: Creating a struct to save to a database where the `ID` is unknown (it stays `0`).
    2. **Receiving Data**: Using a struct field as a placeholder for a database return value (like a `RETURNING id` or `LastInsertId`).
- **`omitempty` tag**: Used in JSON tags (`` `json:"id,omitempty"` ``) to tell Go to skip that field during JSON encoding if it holds its zero value.
