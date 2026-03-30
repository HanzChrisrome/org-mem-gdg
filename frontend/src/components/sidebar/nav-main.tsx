import {
  SidebarGroup,
  SidebarGroupContent,
  SidebarGroupLabel,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "@/components/ui/sidebar";
import { Link, useLocation } from "react-router-dom";

export function NavMain({
  items,
}: {
  items: {
    title: string;
    url: string;
    icon?: React.ReactNode;
    activeMatch?: "exact" | "prefix";
  }[];
}) {
  const { pathname } = useLocation();

  const isItemActive = (item: {
    url: string;
    activeMatch?: "exact" | "prefix";
  }) => {
    if (item.activeMatch === "prefix") {
      return pathname === item.url || pathname.startsWith(`${item.url}/`);
    }

    return pathname === item.url;
  };

  return (
    <SidebarGroup>
      <SidebarGroupLabel>Dashboard</SidebarGroupLabel>
      <SidebarGroupContent className="flex flex-col gap-2">
        <SidebarMenu>
          {items.map((item) => (
            <SidebarMenuItem key={item.title}>
              <SidebarMenuButton asChild isActive={isItemActive(item)}>
                <Link to={item.url}>
                  <span className="flex items-center gap-2 w-full">
                    {item.icon}
                    <span>{item.title}</span>
                  </span>
                </Link>
              </SidebarMenuButton>
            </SidebarMenuItem>
          ))}
        </SidebarMenu>
      </SidebarGroupContent>
    </SidebarGroup>
  );
}
