# Go Learning Repository

Welcome to my personal Go learning repository! 

This repository serves as a knowledge base and practical playground for exploring and mastering various concepts, libraries, and best practices in the Go programming language. It is organized into topical modules, each containing documented code examples, guidelines, and automation tools designed to solidify understanding.

## Table of Contents

- [1. Error Handling (`/errors`)](#1-error-handling)

---

## 1. Error Handling
**Path:** [`/errors`](./errors)

A comprehensive guide and set of interactive examples demonstrating how to handle errors idiomatically in Go, aligning with the "errors are values" philosophy, the Uber Go Style Guide, and Domain-Driven Design (DDD) boundaries.

### Key Concepts Covered:
- **Raw Ingredients:** Using `errors.New` vs `fmt.Errorf`.
- **The "Go-nion":** Wrapping and adding context with `%w` vs obfuscating with `%v`.
- **Peeling & Grouping:** Using `errors.Is`, `errors.As`, `errors.AsType` (Go 1.26+), and `errors.Join`.
- **Best Practices:** The "Handle or Return, Never Both" rule, graceful degradation, and the "Don't Panic" principle.
- **DDD Error Boundaries:** The "Split-Brain" pattern, catching low-level infrastructure errors and mapping them into safe, sanitized Domain Errors for public APIs to prevent leaking the "kitchen core".

### Tools & Automation Included:
- **`go-error-style` Skill:** An Agent Skill specification (`SKILL.md`) that enforces error handling codestyles.
- **`go-error-refactor` Subagent:** A configured Claude Code subagent (`.claude/agents`) that autonomously scans directories for double-logging, opaque returns, and missing domain boundaries, and automatically refactors them to standard.
