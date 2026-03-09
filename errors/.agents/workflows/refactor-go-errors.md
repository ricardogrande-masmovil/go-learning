---
description: Automatically discover and refactor Go files to adhere to the project's error handling codestyle.
---

This workflow automates the process of finding non-compliant Go error handling and systematically refactoring it using the `go-error-style` skill.

1. **Discovery**: Ask the user for the parent path where you should start reading code. Once provided, actively scan that specific directory to find non-compliant code. You should use a combination of tools (like `grep_search` or writing a temporary Go AST parsing script) to identify the following patterns within that path:
    - Code that contains a logging statement (e.g., `log.`, `zap.`, `logrus.`) followed by returning an error shortly after (double-logging anti-pattern).
    - Sentinel errors defined in domain layers that should be converted into custom error types (`structs`) to hold more semantic information.
    - Places where errors are returned opaquely (e.g., `fmt.Errorf("db fail")` or `errors.New("fail")`) without wrapping the underlying root cause when the error chain should be preserved.
    - Occurrences of the `"failed to"` string wrapping anti-pattern.
2. **Backlog Creation**: If no non-compliant code is found during the discovery step, notify the user that the code is already compliant and finish the workflow. Otherwise, create or update the `task.md` artifact to list every file discovered in step 1 that requires refactoring.
3. **Execution**: For each file in the backlog:
    - Read the file to understand the context.
    - Apply the `go-error-style` skill to fix the redundant wrapping or raw error leaks.
    - Verify the file compiles by running `go build`.
    - Mark the file as complete `[x]` in `task.md` before moving to the next.
4. **Conclusion**: Once the backlog is empty, summarize the changes made to the user.
