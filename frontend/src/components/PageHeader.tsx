type PageHeaderProps = {
  total: number;
  completed: number;
};

export function PageHeader({ total, completed }: PageHeaderProps) {
  return (
    <header className="hero-card">
      <div>
        <p className="eyebrow">Agent Harness Demo Baseline</p>
        <h1>Todo workspace</h1>
        <p className="hero-copy">
          A lightweight live demo project with CRUD, SQLite persistence, priority, filtering, and
          a clean single-page workflow.
        </p>
      </div>

      <div className="hero-metrics" aria-label="Todo summary">
        <div className="metric-card">
          <span className="metric-label">Total todos</span>
          <strong>{total}</strong>
        </div>
        <div className="metric-card">
          <span className="metric-label">Completed</span>
          <strong>{completed}</strong>
        </div>
      </div>
    </header>
  );
}
