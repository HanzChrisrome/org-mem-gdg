import DashboardLayout from "@/components/layout/DashboardLayout";

export default function ExecutiveDashboard() {
  return (
    <DashboardLayout>
      <h1 className="text-2xl font-bold mb-6">Welcome back, Executive</h1>

      <div className="grid grid-cols-3 gap-4">
        <div className="bg-white p-4 rounded shadow">Total Members</div>

        <div className="bg-white p-4 rounded shadow">Pending Payments</div>

        <div className="bg-white p-4 rounded shadow">Upcoming Events</div>
      </div>
    </DashboardLayout>
  );
}
