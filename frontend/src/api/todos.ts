import { apiRequest } from "./client";
import type { CreateTodoInput, Todo, UpdateTodoInput } from "../features/todos/types";

export function fetchTodos() {
  return apiRequest<Todo[]>("/api/todos");
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
