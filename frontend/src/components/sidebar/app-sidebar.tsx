import * as React from "react";

import { NavMain } from "@/components/sidebar/nav-main";
import { NavMembers } from "@/components/sidebar/nav-members";
import { NavOthers } from "@/components/sidebar/nav-others";
import { NavPayments } from "@/components/sidebar/nav-payments";
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
} from "@/components/ui/sidebar";
import {
  BarChartIcon,
  BuildingIcon,
  CheckCircleIcon,
  CogIcon,
  LayoutDashboardIcon,
  LogsIcon,
  UserPlusIcon,
} from "lucide-react";

import logo from "@/assets/gdgoc-logo.png";
import { NavUser } from "./nav-user";

const data = {
  user: {
    name: "Iorilovesme",
    email: "m@example.com",
    avatar: "/avatars/shadcn.jpg",
  },
  navMain: [
    {
      title: "Overview",
      url: "/dashboard",
      icon: <LayoutDashboardIcon />,
      activeMatch: "exact" as const,
    },
  ],
  members: [
    {
      name: "Member List",
      url: "/members",
      icon: <BarChartIcon />,
      activeMatch: "exact" as const,
    },
    {
      name: "Add New Member",
      url: "/members/new",
      icon: <UserPlusIcon />,
      activeMatch: "exact" as const,
    },
  ],
  payments: [
    {
      name: "Payment Approval",
      url: "/payments/approvals",
      icon: <CheckCircleIcon />,
      activeMatch: "exact" as const,
    },
  ],
  others: [
    {
      name: "Audit Logs",
      url: "#",
      icon: <LogsIcon />,
    },
    {
      name: "Settings",
      url: "#",
      icon: <CogIcon />,
    },
    {
      name: "Organization Settings",
      url: "#",
      icon: <BuildingIcon />,
    },
  ],
};

export function AppSidebar({ ...props }: React.ComponentProps<typeof Sidebar>) {
  return (
    <Sidebar collapsible="offcanvas" {...props}>
      <SidebarHeader>
        <div className="flex items-center justify-center">
          <img src={logo} alt="GDG on Campus Logo" className="h-14 w-auto" />
        </div>
      </SidebarHeader>
      <SidebarContent>
        <NavMain items={data.navMain} />
        <NavMembers items={data.members} />
        <NavPayments items={data.payments} />
        <NavOthers items={data.others} />
      </SidebarContent>
      <SidebarFooter>
        <NavUser user={data.user} />
      </SidebarFooter>
    </Sidebar>
  );
}
