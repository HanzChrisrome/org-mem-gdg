import { ChartAreaInteractive } from "@/components/chart-area-interactive";
import { DataTable } from "@/components/data-table";
import DashboardLayout from "@/components/layout/DashboardLayout";
import { SectionCards } from "@/components/section-cards";

import data from "./data.json";

export default function ExecutiveDashboard() {
  return (
    <DashboardLayout>
      <section className="flex flex-col gap-4 py-4 md:gap-6 md:py-6">
        <h1 className="text-2xl font-semibold tracking-tight px-4 lg:px-6">
          Here's an overview of your organization activity.
        </h1>
        <SectionCards />
        <div className="px-4 lg:px-6">
          <ChartAreaInteractive />
        </div>
        <DataTable data={data} />
      </section>
    </DashboardLayout>
  );
}
