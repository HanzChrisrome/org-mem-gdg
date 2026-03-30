import {
  Breadcrumb,
  BreadcrumbItem,
  BreadcrumbLink,
  BreadcrumbList,
  BreadcrumbPage,
  BreadcrumbSeparator,
} from "@/components/ui/breadcrumb";
import { Separator } from "@/components/ui/separator";
import { SidebarTrigger } from "@/components/ui/sidebar";
import { Fragment } from "react";
import { Link, useLocation } from "react-router-dom";

type Crumb = {
  label: string;
  to?: string;
  hideOnMobile?: boolean;
};

function getBreadcrumbs(pathname: string): Crumb[] {
  if (pathname === "/dashboard") {
    return [{ label: "Overview" }];
  }

  if (pathname === "/members") {
    return [
      { label: "Overview", to: "/dashboard", hideOnMobile: true },
      { label: "Member List" },
    ];
  }

  if (pathname === "/members/new") {
    return [
      { label: "Overview", to: "/dashboard", hideOnMobile: true },
      { label: "Member List", to: "/members", hideOnMobile: true },
      { label: "Add New Member" },
    ];
  }

  if (pathname === "/payments/approvals") {
    return [
      { label: "Overview", to: "/dashboard", hideOnMobile: true },
      { label: "Payment Approval" },
    ];
  }

  return [{ label: "Dashboard" }];
}

export function SiteHeader() {
  const { pathname } = useLocation();
  const breadcrumbs = getBreadcrumbs(pathname);

  return (
    <header className="flex h-(--header-height) shrink-0 items-center gap-2 border-b transition-[width,height] ease-linear group-has-data-[collapsible=icon]/sidebar-wrapper:h-(--header-height)">
      <div className="flex w-full items-center gap-1 px-4 lg:gap-2 lg:px-6">
        <SidebarTrigger className="-ml-1" />
        <Separator orientation="vertical" className="mx-2" />
        <Breadcrumb>
          <BreadcrumbList>
            {breadcrumbs.map((crumb, index) => {
              const isLast = index === breadcrumbs.length - 1;
              const visibilityClass = crumb.hideOnMobile
                ? "hidden sm:inline-flex"
                : "inline-flex";

              return (
                <Fragment key={`${crumb.label}-${index}`}>
                  <BreadcrumbItem className={visibilityClass}>
                    {isLast || !crumb.to ? (
                      <BreadcrumbPage>{crumb.label}</BreadcrumbPage>
                    ) : (
                      <BreadcrumbLink asChild>
                        <Link to={crumb.to}>{crumb.label}</Link>
                      </BreadcrumbLink>
                    )}
                  </BreadcrumbItem>
                  {!isLast && (
                    <BreadcrumbSeparator
                      className={
                        crumb.hideOnMobile
                          ? "hidden sm:inline-flex"
                          : "inline-flex"
                      }
                    />
                  )}
                </Fragment>
              );
            })}
          </BreadcrumbList>
        </Breadcrumb>
      </div>
    </header>
  );
}
