import type { FormEvent } from "react";
import { useEffect, useState } from "react";

import { createTodo, deleteTodo, fetchTodos, updateTodo } from "../../api/todos";
import { PageHeader } from "../../components/PageHeader";
import { TodoCard } from "./TodoCard";
import { TodoComposer } from "./TodoComposer";
import type { Todo, TodoDraft } from "./types";

const emptyDraft: TodoDraft = {
  title: "",
  notes: ""
};

export function TodoPage() {
  const [todos, setTodos] = useState<Todo[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isCreating, setIsCreating] = useState(false);
  const [activeMutationId, setActiveMutationId] = useState<number | null>(null);
  const [draft, setDraft] = useState<TodoDraft>(emptyDraft);
  const [editingId, setEditingId] = useState<number | null>(null);
  const [editDraft, setEditDraft] = useState<TodoDraft>(emptyDraft);

  useEffect(() => {
    void loadTodos();
  }, []);

  async function loadTodos() {
    setIsLoading(true);
    setError(null);

    try {
      const items = await fetchTodos();
      setTodos(items);
    } catch (loadError) {
      setError(getErrorMessage(loadError, "Failed to load todos."));
    } finally {
      setIsLoading(false);
    }
  }

  function updateDraft(field: keyof TodoDraft, value: string) {
    setDraft((current) => ({
      ...current,
      [field]: value
    }));
  }

  function updateEditDraft(field: keyof TodoDraft, value: string) {
    setEditDraft((current) => ({
      ...current,
      [field]: value
    }));
  }

  async function handleCreate(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const title = draft.title.trim();

    if (!title) {
      setError("Title is required.");
      return;
    }

    setIsCreating(true);
    setError(null);

    try {
      const created = await createTodo({
        title,
        notes: draft.notes.trim()
      });
      setTodos((current) => [created, ...current]);
      setDraft(emptyDraft);
    } catch (createError) {
      setError(getErrorMessage(createError, "Failed to create todo."));
    } finally {
      setIsCreating(false);
    }
  }

  function startEditing(todo: Todo) {
    setEditingId(todo.id);
    setEditDraft({
      title: todo.title,
      notes: todo.notes
    });
    setError(null);
  }

  function cancelEditing() {
    setEditingId(null);
    setEditDraft(emptyDraft);
  }

  async function saveEdit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (editingId === null) {
      return;
    }

    const title = editDraft.title.trim();
    if (!title) {
      setError("Title is required.");
      return;
    }

    setActiveMutationId(editingId);
    setError(null);

    try {
      const updated = await updateTodo(editingId, {
        title,
        notes: editDraft.notes.trim()
      });
      setTodos((current) => current.map((todo) => (todo.id === editingId ? updated : todo)));
      cancelEditing();
    } catch (saveError) {
      setError(getErrorMessage(saveError, "Failed to save changes."));
    } finally {
      setActiveMutationId(null);
    }
  }

  async function toggleTodo(todo: Todo) {
    setActiveMutationId(todo.id);
    setError(null);

    try {
      const updated = await updateTodo(todo.id, {
        completed: !todo.completed
      });
      setTodos((current) => current.map((item) => (item.id === todo.id ? updated : item)));
    } catch (toggleError) {
      setError(getErrorMessage(toggleError, "Failed to update the todo."));
    } finally {
      setActiveMutationId(null);
    }
  }

  async function removeTodo(id: number) {
    setActiveMutationId(id);
    setError(null);

    try {
      await deleteTodo(id);
      setTodos((current) => current.filter((todo) => todo.id !== id));
      if (editingId === id) {
        cancelEditing();
      }
    } catch (deleteError) {
      setError(getErrorMessage(deleteError, "Failed to delete the todo."));
    } finally {
      setActiveMutationId(null);
    }
  }

  const completedCount = todos.filter((todo) => todo.completed).length;

  return (
    <main className="page-shell">
      <div className="page-inner">
        <PageHeader total={todos.length} completed={completedCount} />

        {error ? (
          <div className="alert-banner" role="alert">
            {error}
          </div>
        ) : null}

        {(isCreating || activeMutationId !== null) && <p className="mutating-copy">Saving changes...</p>}

        <TodoComposer draft={draft} isSaving={isCreating} onChange={updateDraft} onSubmit={handleCreate} />

        <section className="panel-card">
          <div className="section-heading">
            <div>
              <p className="eyebrow">Todos</p>
              <h2>Current list</h2>
            </div>
            <p className="section-copy">The baseline intentionally has no priority or filtering controls yet.</p>
          </div>

          {isLoading ? (
            <div className="state-card">Loading todos...</div>
          ) : todos.length === 0 ? (
            <div className="state-card">
              <h3>No todos yet</h3>
              <p>Create your first todo above.</p>
            </div>
          ) : (
            <div className="todo-grid">
              {todos.map((todo) => (
                <TodoCard
                  key={todo.id}
                  todo={todo}
                  isEditing={editingId === todo.id}
                  isSaving={activeMutationId === todo.id}
                  editDraft={editingId === todo.id ? editDraft : { title: todo.title, notes: todo.notes }}
                  onEditStart={() => startEditing(todo)}
                  onEditChange={updateEditDraft}
                  onEditCancel={cancelEditing}
                  onEditSubmit={saveEdit}
                  onToggle={() => toggleTodo(todo)}
                  onDelete={() => removeTodo(todo.id)}
                />
              ))}
            </div>
          )}
        </section>
      </div>
    </main>
  );
}

function getErrorMessage(error: unknown, fallback: string) {
  if (error instanceof Error && error.message) {
    return error.message;
  }

  return fallback;
}
