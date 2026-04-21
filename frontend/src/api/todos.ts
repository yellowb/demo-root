import { apiRequest } from "./client";
import type { CreateTodoInput, Todo, TodoListFilters, UpdateTodoInput } from "../features/todos/types";

export function fetchTodos(filters: TodoListFilters = {}) {
  const params = new URLSearchParams();
  if (filters.completed !== undefined) {
    params.set("completed", String(filters.completed));
  }
  if (filters.priority) {
    params.set("priority", filters.priority);
  }

  const query = params.toString();
  return apiRequest<Todo[]>(query ? `/api/todos?${query}` : "/api/todos");
}

export function createTodo(input: CreateTodoInput) {
  return apiRequest<Todo>("/api/todos", {
    method: "POST",
    body: JSON.stringify(input)
  });
}

export function updateTodo(id: number, input: UpdateTodoInput) {
  return apiRequest<Todo>(`/api/todos/${id}`, {
    method: "PATCH",
    body: JSON.stringify(input)
  });
}

export function deleteTodo(id: number) {
  return apiRequest<void>(`/api/todos/${id}`, {
    method: "DELETE"
  });
}
