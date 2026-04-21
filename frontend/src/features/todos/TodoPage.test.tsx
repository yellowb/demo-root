import { render, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import App from "../../App";

type Todo = {
  id: number;
  title: string;
  notes: string;
  completed: boolean;
  priority: "low" | "medium" | "high";
  created_at: string;
  updated_at: string;
};

const baseTodoList: Todo[] = [
  {
    id: 1,
    title: "Prepare Agent Harness 分享提纲",
    notes: "补齐 live demo 的上下文和 framing",
    completed: false,
    priority: "high",
    created_at: "2026-04-18T09:00:00Z",
    updated_at: "2026-04-18T09:00:00Z"
  },
  {
    id: 2,
    title: "补充后端测试",
    notes: "覆盖 CRUD 和 handler 主路径",
    completed: true,
    priority: "low",
    created_at: "2026-04-18T10:00:00Z",
    updated_at: "2026-04-18T10:00:00Z"
  }
];

describe("TodoPage", () => {
  beforeEach(() => {
    installFetchServer(baseTodoList);
  });

  afterEach(() => {
    vi.restoreAllMocks();
  });

  it("renders todos, creates a todo, edits it, toggles completion, and deletes it", async () => {
    const user = userEvent.setup();

    render(<App />);

    expect(screen.getByText(/Loading todos/i)).toBeInTheDocument();

    await screen.findByText("补充后端测试");
    expect(screen.getByText("Prepare Agent Harness 分享提纲")).toBeInTheDocument();
    expect(within(screen.getByTestId("todo-card-1")).getByText(/High priority/i)).toBeInTheDocument();
    expect(within(screen.getByTestId("todo-card-2")).getByText(/Low priority/i)).toBeInTheDocument();

    await user.type(screen.getByLabelText(/Title/i), "Draft baseline README");
    await user.type(screen.getByLabelText(/Notes/i), "Capture setup, test, and reset commands");
    await user.selectOptions(screen.getByLabelText(/^Priority$/i), "high");
    await user.click(screen.getByRole("button", { name: /Add todo/i }));

    await screen.findByText("Draft baseline README");
    expect(within(screen.getByTestId("todo-card-3")).getByText(/High priority/i)).toBeInTheDocument();

    const newTodoCard = screen.getByTestId("todo-card-3");
    await user.click(within(newTodoCard).getByRole("button", { name: /Edit/i }));

    const editTitle = within(newTodoCard).getByLabelText(/Edit title/i);
    await user.clear(editTitle);
    await user.type(editTitle, "Draft baseline README v2");

    const editNotes = within(newTodoCard).getByLabelText(/Edit notes/i);
    await user.clear(editNotes);
    await user.type(editNotes, "Call out that priority is intentionally absent");
    await user.selectOptions(within(newTodoCard).getByLabelText(/Edit priority/i), "low");
    await user.click(within(newTodoCard).getByRole("button", { name: /Save changes/i }));

    await screen.findByText("Draft baseline README v2");
    expect(within(newTodoCard).getByText(/Low priority/i)).toBeInTheDocument();

    await user.click(within(newTodoCard).getByRole("button", { name: /Mark complete/i }));
    await screen.findAllByText(/Completed/i);

    await user.click(within(newTodoCard).getByRole("button", { name: /Delete/i }));

    await waitFor(() => {
      expect(screen.queryByText("Draft baseline README v2")).not.toBeInTheDocument();
    });
  });

  it("shows the empty state when no todos are returned", async () => {
    installFetchServer([]);

    render(<App />);

    await screen.findByText(/No todos yet/i);
    expect(screen.getByText(/Create your first todo above/i)).toBeInTheDocument();
  });

  it("filters todos by completion status and priority", async () => {
    const user = userEvent.setup();

    render(<App />);

    await screen.findByText("补充后端测试");

    await user.selectOptions(screen.getByLabelText(/Completion filter/i), "active");
    await user.selectOptions(screen.getByLabelText(/Priority filter/i), "high");

    await waitFor(() => {
      expect(screen.getByText("Prepare Agent Harness 分享提纲")).toBeInTheDocument();
      expect(screen.queryByText("补充后端测试")).not.toBeInTheDocument();
    });
  });
});

function installFetchServer(initialTodos: Todo[]) {
  let todos = [...initialTodos];

  vi.stubGlobal(
    "fetch",
    vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      const url = typeof input === "string" ? input : input.toString();
      const method = init?.method ?? "GET";

      if (url.startsWith("/api/todos") && method === "GET") {
        const parsedURL = new URL(url, "http://localhost");
        const completed = parsedURL.searchParams.get("completed");
        const priority = parsedURL.searchParams.get("priority");
        const filtered = todos.filter((todo) => {
          const completedMatches = completed === null || String(todo.completed) === completed;
          const priorityMatches = priority === null || todo.priority === priority;
          return completedMatches && priorityMatches;
        });

        return jsonResponse(filtered);
      }

      if (url === "/api/todos" && method === "POST") {
        const payload = JSON.parse(String(init?.body ?? "{}")) as {
          title: string;
          notes: string;
          priority?: Todo["priority"];
        };
        const nextId = todos.reduce((max, todo) => Math.max(max, todo.id), 0) + 1;
        const todo: Todo = {
          id: nextId,
          title: payload.title,
          notes: payload.notes,
          completed: false,
          priority: payload.priority ?? "medium",
          created_at: "2026-04-19T09:00:00Z",
          updated_at: "2026-04-19T09:00:00Z"
        };
        todos = [todo, ...todos];
        return jsonResponse(todo, 201);
      }

      const match = url.match(/^\/api\/todos\/(\d+)$/);
      if (!match) {
        throw new Error(`Unexpected request: ${method} ${url}`);
      }

      const todoId = Number(match[1]);

      if (method === "PATCH") {
        const payload = JSON.parse(String(init?.body ?? "{}")) as Partial<Todo>;
        const current = todos.find((todo) => todo.id === todoId);
        if (!current) {
          return jsonResponse({ error: "todo not found" }, 404);
        }

        const updated: Todo = {
          ...current,
          title: payload.title ?? current.title,
          notes: payload.notes ?? current.notes,
          completed: payload.completed ?? current.completed,
          priority: payload.priority ?? current.priority,
          updated_at: "2026-04-19T09:30:00Z"
        };
        todos = todos.map((todo) => (todo.id === todoId ? updated : todo));
        return jsonResponse(updated);
      }

      if (method === "DELETE") {
        todos = todos.filter((todo) => todo.id !== todoId);
        return new Response(null, { status: 204 });
      }

      throw new Error(`Unexpected request: ${method} ${url}`);
    })
  );
}

function jsonResponse(payload: unknown, status = 200) {
  return new Response(JSON.stringify(payload), {
    status,
    headers: {
      "Content-Type": "application/json"
    }
  });
}
