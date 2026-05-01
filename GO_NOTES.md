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
