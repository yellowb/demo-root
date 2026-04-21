# Live Demo Context

## Technical Talk Goal

This repository belongs to an Agent Harness technical talk.

The talk is not a tool review and not a "watch AI write code" show. Its goal is to help developers understand this capability path:

```text
Prompt -> Context -> Harness
```

The key message is that real coding tasks need more than a clear prompt. They also need the right repository context, executable tools, validation commands, and a feedback loop.

## Harness Engineering Points

This demo should make four repository-level Harness practices visible:

1. Context engineering: `AGENTS.md` defines project positioning, stack, boundaries, API shape, and workflow constraints; `openspec/specs/todo-management/spec.md` captures the current Todo baseline as an executable spec reference.
2. Unified execution entrypoints: `Makefile` exposes `setup`, `dev`, `test`, and `reset-db` so the Agent does not need to guess how to run the project.
3. Validation loop: `make test` validates OpenSpec specs, backend checks, frontend checks, and Codex hook configuration; the Codex `Stop` hook runs it before a final response when repository changes exist.
4. Specified change entrypoint: OpenSpec lets the Agent turn a feature request into `proposal.md`, `design.md`, and `tasks.md` before editing application code.

## Case 1 vs Case 2

- Case 1 uses Claude Code to explain Harness mechanisms: context organization, tool use, permission control, sandboxing, MCP, hooks, and subagents.
- Case 2 uses this Todo List project to show how Harness engineering lands in a web developer's daily workflow.

Case 2 should not re-explain Harness terminology or become an OpenSpec tutorial. It should show how a Coding Agent handles a real repository task with context, specifications, tools, and verification already prepared.

## What This Demo Should Show

The live demo should help the audience observe this workflow:

```text
task description -> inspect repo/spec -> create OpenSpec change -> update frontend/backend/persistence -> run checks -> review result -> converge
```

The point is not that the AI can add a feature. The point is that a well-prepared repository gives the Agent a clearer, safer, and more verifiable working environment.

## Baseline State

This repository is the baseline before the live demo.

- It is already runnable.
- It has demo seed data.
- It supports Todo CRUD and SQLite persistence.
- It intentionally does not include `priority + filtering` yet.

Do not add priority, filtering, search, tags, due dates, auth, or multi-user collaboration unless explicitly asked.

## Recommended Live Demo Prompt

Use this prompt in a fresh Codex session opened at this repository root:

```text
Please add priority and filtering support to this Todo List app.

Requirements:
1. Add a priority field to todos with three values: low, medium, high.
2. Allow creating and editing priority from the UI.
3. Show each todo's priority clearly in the list.
4. Add filtering support so I can filter by completion status and priority.
5. Update the backend API and persistence layer as needed.
6. Update tests and basic validation commands.
7. Update seed data so the result is visually obvious in the demo.

Before editing application code:
1. Inspect the repo context and the current OpenSpec baseline spec.
2. Use OpenSpec propose to create an `add-priority-filtering` change with proposal, design, specs, and tasks.
3. Briefly explain the plan from those artifacts.

After changes, run the relevant verification commands and summarize what changed.
```

## Demo Success Criteria

- The Agent reads the repository context instead of relying only on the prompt.
- The Agent uses OpenSpec artifacts to explain the change scope, implementation tasks, and acceptance basis before editing code.
- The Agent identifies that the feature crosses backend schema, repository logic, API handlers, frontend API client, UI, seed data, and tests.
- The validation loop is visible through `make test`, including OpenSpec spec validation.
- The final result is explainable as a Harness workflow: context, tools, validation, and feedback.
