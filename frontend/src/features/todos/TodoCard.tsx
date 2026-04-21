import type { FormEvent } from "react";

import { TODO_PRIORITIES, type Todo, type TodoDraft } from "./types";

type TodoCardProps = {
  todo: Todo;
  isEditing: boolean;
  isSaving: boolean;
  editDraft: TodoDraft;
  onEditStart: () => void;
  onEditChange: (field: keyof TodoDraft, value: string) => void;
  onEditCancel: () => void;
  onEditSubmit: (event: FormEvent<HTMLFormElement>) => void;
  onToggle: () => void;
  onDelete: () => void;
};

export function TodoCard({
  todo,
  isEditing,
  isSaving,
  editDraft,
  onEditStart,
  onEditChange,
  onEditCancel,
  onEditSubmit,
  onToggle,
  onDelete
}: TodoCardProps) {
  return (
    <article className="todo-card" data-testid={`todo-card-${todo.id}`}>
      <div className="todo-card-top">
        <div>
          <p className="todo-meta">{formatTimestamp(todo.updated_at)}</p>
          <h3>{todo.title}</h3>
        </div>
        <div className="todo-pill-group">
          <span className={`priority-pill priority-pill-${todo.priority}`}>{formatPriority(todo.priority)}</span>
          <span className={todo.completed ? "status-pill status-pill-done" : "status-pill"}>
            {todo.completed ? "Completed" : "In progress"}
          </span>
        </div>
      </div>

      {isEditing ? (
        <form className="todo-form inline-form" onSubmit={onEditSubmit}>
          <label className="field">
            <span>Edit title</span>
            <input
              name="edit-title"
              type="text"
              value={editDraft.title}
              onChange={(event) => onEditChange("title", event.target.value)}
              disabled={isSaving}
            />
          </label>

          <label className="field">
            <span>Edit notes</span>
            <textarea
              name="edit-notes"
              value={editDraft.notes}
              onChange={(event) => onEditChange("notes", event.target.value)}
              rows={3}
              disabled={isSaving}
            />
          </label>

          <label className="field">
            <span>Edit priority</span>
            <select
              name="edit-priority"
              value={editDraft.priority}
              onChange={(event) => onEditChange("priority", event.target.value)}
              disabled={isSaving}
            >
              {TODO_PRIORITIES.map((priority) => (
                <option key={priority} value={priority}>
                  {formatPriority(priority)}
                </option>
              ))}
            </select>
          </label>

          <div className="inline-actions">
            <button className="primary-button" type="submit" disabled={isSaving}>
              {isSaving ? "Saving..." : "Save changes"}
            </button>
            <button className="ghost-button" type="button" onClick={onEditCancel} disabled={isSaving}>
              Cancel
            </button>
          </div>
        </form>
      ) : (
        <>
          <p className="todo-notes">{todo.notes || "No notes yet."}</p>
          <div className="todo-actions">
            <button className="ghost-button" type="button" onClick={onEditStart} disabled={isSaving}>
              Edit
            </button>
            <button className="ghost-button" type="button" onClick={onToggle} disabled={isSaving}>
              {todo.completed ? "Mark active" : "Mark complete"}
            </button>
            <button className="danger-button" type="button" onClick={onDelete} disabled={isSaving}>
              Delete
            </button>
          </div>
        </>
      )}
    </article>
  );
}

function formatTimestamp(value: string) {
  return new Intl.DateTimeFormat("en-US", {
    dateStyle: "medium",
    timeStyle: "short"
  }).format(new Date(value));
}

function formatPriority(value: string) {
  return `${value.charAt(0).toUpperCase()}${value.slice(1)} priority`;
}
