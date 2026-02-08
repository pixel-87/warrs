---
on: 
  push:
    paths:
      - '**/*.go'
permissions: read-all
safe-outputs:
  create-pull-request:
    labels:
      - "ai-generated"
      - "testing"
---
# Universal Pure-Logic Tester

Analyze all `.go` files in this repository. 
I want to increase test coverage using **Pure Functional Testing** principles.

**Goal:** Find logic that is currently untested or undertested and write **Table-Driven Tests** for it without using mocks.

**Strategy:**
1. **Identify Candidates:** Scan the codebase for functions that:
   - Take simple inputs (strings, structs, slices, primitives).
   - Return values or errors.
   - Do NOT directly call `http`, `sql`, or `os` packages (pure logic).
   
2. **Generate Tests:** For every candidate found:
   - Check if a `_test.go` file exists. If not, create it.
   - Write a **Table-Driven Test** structure.
   - **Constraint:** DO NOT use mocks. If a function is hard to test without mocks, skip it.

3. **Robustness Check (The "Fuzz" Factor):**
   - When populating the test cases, do not just use "happy path" data.
   - Include **Edge Cases**: Empty strings, nil slices, negative integers, huge payloads, malformed JSON/XML.
   - We want to ensure the function handles garbage input gracefully (returns error) rather than panicking.

4. **Execution:**
   - Run the newly generated tests.
   - If they fail, adjust the *test expectation* (assuming the code is right) OR fix the code if it's a clear bug (like a nil pointer panic).

5. **Output:**
   - Create a Pull Request with the new tests.
   - **Title the PR:** "[cai-tests] New pure logic tests"
   - Group multiple test files into a single PR if possible.
