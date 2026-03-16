import { Link } from "react-router-dom";

export default function Sidebar() {
  return (
    <aside className="w-64 bg-white border-r p-4">
      <h2 className="text-xl font-bold mb-6">GDG Admin</h2>

      <nav className="space-y-3">
        <Link to="/dashboard" className="block">
          Dashboard
        </Link>

        <Link to="/members" className="block">
          Members
        </Link>

        <Link to="/payments" className="block">
          Payments
        </Link>

        <Link to="/reports" className="block">
          Reports
        </Link>

        <Link to="/audit" className="block">
          Audit Logs
        </Link>
      </nav>
    </aside>
  );
}
