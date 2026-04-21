export const TODO_PRIORITIES = ["low", "medium", "high"] as const;

export type TodoPriority = (typeof TODO_PRIORITIES)[number];

export type Todo = {
  id: number;
  title: string;
  notes: string;
  completed: boolean;
  priority: TodoPriority;
  created_at: string;
  updated_at: string;
};

export type CreateTodoInput = {
  title: string;
  notes: string;
  priority: TodoPriority;
};

export type UpdateTodoInput = {
  title?: string;
  notes?: string;
  completed?: boolean;
  priority?: TodoPriority;
};

export type TodoListFilters = {
  completed?: boolean;
  priority?: TodoPriority;
};

export type TodoDraft = {
  title: string;
  notes: string;
  priority: TodoPriority;
};
