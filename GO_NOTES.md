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
