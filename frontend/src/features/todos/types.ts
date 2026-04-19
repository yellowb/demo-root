export type Todo = {
  id: number;
  title: string;
  notes: string;
  completed: boolean;
  created_at: string;
  updated_at: string;
};

export type CreateTodoInput = {
  title: string;
  notes: string;
};

export type UpdateTodoInput = {
  title?: string;
  notes?: string;
  completed?: boolean;
};

export type TodoDraft = {
  title: string;
  notes: string;
};
