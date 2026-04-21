import type { FormEvent } from "react";

import { TODO_PRIORITIES, type TodoDraft } from "./types";

type TodoComposerProps = {
  draft: TodoDraft;
  isSaving: boolean;
  onChange: (field: keyof TodoDraft, value: string) => void;
  onSubmit: (event: FormEvent<HTMLFormElement>) => void;
};

export function TodoComposer({ draft, isSaving, onChange, onSubmit }: TodoComposerProps) {
  return (
    <section className="panel-card">
      <div className="section-heading">
        <div>
          <p className="eyebrow">Create</p>
          <h2>Add a new todo</h2>
        </div>
        <p className="section-copy">Keep the form at the top so new tasks are always one step away.</p>
      </div>

      <form className="todo-form" onSubmit={onSubmit}>
        <label className="field">
          <span>Title</span>
          <input
            name="title"
            type="text"
            value={draft.title}
            onChange={(event) => onChange("title", event.target.value)}
            placeholder="Draft live demo checklist"
            disabled={isSaving}
          />
        </label>

        <label className="field">
          <span>Notes</span>
          <textarea
            name="notes"
            value={draft.notes}
            onChange={(event) => onChange("notes", event.target.value)}
            rows={3}
            placeholder="Add a small amount of context for later editing"
            disabled={isSaving}
          />
        </label>

        <label className="field">
          <span>Priority</span>
          <select
            name="priority"
            value={draft.priority}
            onChange={(event) => onChange("priority", event.target.value)}
            disabled={isSaving}
          >
            {TODO_PRIORITIES.map((priority) => (
              <option key={priority} value={priority}>
                {formatPriority(priority)}
              </option>
            ))}
          </select>
        </label>

        <div className="form-actions">
          <button className="primary-button" type="submit" disabled={isSaving}>
            {isSaving ? "Saving..." : "Add todo"}
          </button>
        </div>
      </form>
    </section>
  );
}

function formatPriority(value: string) {
  return `${value.charAt(0).toUpperCase()}${value.slice(1)} priority`;
}
