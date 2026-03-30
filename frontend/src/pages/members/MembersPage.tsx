import { listMembers } from "@/api/members";
import DashboardLayout from "@/components/layout/DashboardLayout";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import {
  Table,
  TableBody,
  TableCell,
  TableFooter,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import type { MemberWithPaymentResponse } from "@/types/member";
import { keepPreviousData, useQuery } from "@tanstack/react-query";
import type { AxiosError } from "axios";
import { Loader2 } from "lucide-react";
import { useEffect, useMemo, useState } from "react";

import { Link } from "react-router-dom";
type StatusFilter = "all" | "pending" | "active" | "inactive" | "resubmitted";
type PaymentFilter =
  | "all"
  | "pending"
  | "approved"
  | "rejected"
  | "resubmitted";

const MIN_LOADING_MS = 2000;

function wait(ms: number) {
  return new Promise<void>((resolve) => {
    window.setTimeout(resolve, ms);
  });
}

function normalizeStatus(status?: string | null) {
  return (status ?? "").trim().toLowerCase();
}

function titleCaseStatus(status?: string | null) {
  const normalized = normalizeStatus(status);
  if (!normalized) {
    return "No Record";
  }

  return normalized
    .split(" ")
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
}

function statusBadgeClass(status?: string | null) {
  const normalized = normalizeStatus(status);

  if (normalized === "approved" || normalized === "active") {
    return "bg-emerald-100 text-emerald-700 border-emerald-200 hover:bg-emerald-100";
  }

  if (normalized === "pending") {
    return "bg-amber-100 text-amber-700 border-amber-200 hover:bg-amber-100";
  }

  if (normalized === "resubmitted") {
    return "bg-sky-100 text-sky-700 border-sky-200 hover:bg-sky-100";
  }

  if (!normalized) {
    return "bg-slate-100 text-slate-700 border-slate-200 hover:bg-slate-100";
  }

  return "bg-rose-100 text-rose-700 border-rose-200 hover:bg-rose-100";
}

function formatDateTime(value: string) {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }

  return new Intl.DateTimeFormat("en-PH", {
    year: "numeric",
    month: "short",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
  }).format(date);
}

export default function MembersPage() {
  const [searchInput, setSearchInput] = useState("");
  const [debouncedSearch, setDebouncedSearch] = useState("");
  const [paymentFilter, setPaymentFilter] = useState<PaymentFilter>("all");
  const [registrationFilter, setRegistrationFilter] =
    useState<StatusFilter>("all");
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(10);

  const {
    data: members = [],
    isPending,
    isFetching,
    error,
    refetch,
  } = useQuery<MemberWithPaymentResponse[], AxiosError<{ error?: string }>>({
    queryKey: ["members", { registrationFilter, q: debouncedSearch }],
    queryFn: async () => {
      const [data] = await Promise.all([
        listMembers({
          q: debouncedSearch || undefined,
          status: registrationFilter === "all" ? undefined : registrationFilter,
        }),
        wait(MIN_LOADING_MS),
      ]);

      return data;
    },
    placeholderData: keepPreviousData,
  });

  useEffect(() => {
    const timeoutId = window.setTimeout(() => {
      setDebouncedSearch(searchInput.trim());
    }, 300);

    return () => window.clearTimeout(timeoutId);
  }, [searchInput]);

  const errorMessage = useMemo(() => {
    if (!error) {
      return null;
    }

    if (error.response?.status === 401) {
      return "Session expired or missing. Please log in again to load members.";
    }

    return "Unable to load members. Please try again.";
  }, [error]);

  const filteredMembers = useMemo(() => {
    return members.filter((member) => {
      const paymentStatus = normalizeStatus(member.latest_payment_status);
      const matchesPayment =
        paymentFilter === "all" || paymentStatus === paymentFilter;

      return matchesPayment;
    });
  }, [members, paymentFilter]);

  const totalMembers = filteredMembers.length;
  const totalPages = Math.max(1, Math.ceil(totalMembers / pageSize));
  const safeCurrentPage = Math.min(currentPage, totalPages);
  const startIndex = (safeCurrentPage - 1) * pageSize;
  const paginatedMembers = filteredMembers.slice(
    startIndex,
    startIndex + pageSize,
  );
  const visibleStart = totalMembers === 0 ? 0 : startIndex + 1;
  const visibleEnd = Math.min(startIndex + pageSize, totalMembers);

  return (
    <DashboardLayout>
      <section className="space-y-6">
        <div className="flex flex-col items-start gap-4 lg:flex-row lg:items-center lg:justify-between">
          <div>
            <h1 className="text-xl md:text-2xl lg:text-3xl font-semibold tracking-tight leading-normal">
              Member List
            </h1>
            <p className="text-muted-foreground">
              View and monitor members with their registration and payment
              statuses.
            </p>
          </div>
          <div className="flex items-center gap-2">
            <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
              <div className="flex flex-1 flex-col gap-3 sm:flex-row">
                <Input
                  placeholder="Search by name, student ID, or email"
                  className="sm:max-w-xs"
                  value={searchInput}
                  onChange={(event) => {
                    setCurrentPage(1);
                    setSearchInput(event.target.value);
                  }}
                />
                <Select
                  value={paymentFilter}
                  onValueChange={(value) => {
                    setCurrentPage(1);
                    setPaymentFilter(value as PaymentFilter);
                  }}
                >
                  <SelectTrigger className="w-full sm:w-44">
                    <SelectValue placeholder="Payment status" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All payments</SelectItem>
                    <SelectItem value="pending">Pending</SelectItem>
                    <SelectItem value="approved">Approved</SelectItem>
                    <SelectItem value="rejected">Rejected</SelectItem>
                    <SelectItem value="resubmitted">Resubmitted</SelectItem>
                  </SelectContent>
                </Select>
                <Select
                  value={registrationFilter}
                  onValueChange={(value) => {
                    setCurrentPage(1);
                    setRegistrationFilter(value as StatusFilter);
                  }}
                >
                  <SelectTrigger className="w-full sm:w-44">
                    <SelectValue placeholder="Registration status" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="all">All registrations</SelectItem>
                    <SelectItem value="active">Active</SelectItem>
                    <SelectItem value="pending">Pending</SelectItem>
                    <SelectItem value="inactive">Inactive</SelectItem>
                    <SelectItem value="resubmitted">Resubmitted</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              <Button asChild>
                <Link to="/members/new">Add new member</Link>
              </Button>
            </div>
          </div>
        </div>

        <Card className="p-0">
          <CardContent className="p-0">
            <Table>
              <TableHeader className="bg-slate-100">
                <TableRow>
                  <TableHead>Member</TableHead>
                  <TableHead>Student ID</TableHead>
                  <TableHead>Course</TableHead>
                  <TableHead>Contact</TableHead>
                  <TableHead>Registration</TableHead>
                  <TableHead>Payment</TableHead>
                  <TableHead>Created</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody className="min-w-full">
                {errorMessage ? (
                  <TableRow>
                    <TableCell colSpan={7} className="py-8 text-center">
                      <div className="flex flex-col items-center gap-3">
                        <p className="text-sm text-rose-600">{errorMessage}</p>
                        <Button
                          variant="outline"
                          size="sm"
                          onClick={() => void refetch()}
                        >
                          Retry
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : isPending || (isFetching && totalMembers === 0) ? (
                  <TableRow>
                    <TableCell
                      colSpan={7}
                      className="py-8 text-center text-muted-foreground"
                    >
                      <div className="flex items-center justify-center gap-2">
                        <Loader2 className="h-4 w-4 animate-spin" />
                        <span>Loading members...</span>
                      </div>
                    </TableCell>
                  </TableRow>
                ) : totalMembers === 0 ? (
                  <TableRow>
                    <TableCell
                      colSpan={7}
                      className="py-8 text-center text-muted-foreground"
                    >
                      No members found for the current search and filters.
                    </TableCell>
                  </TableRow>
                ) : (
                  <>
                    {paginatedMembers.map((member) => (
                      <TableRow key={member.member_id}>
                        <TableCell>
                          <div className="flex flex-col">
                            <span className="font-medium">{member.name}</span>
                            <span className="text-xs text-muted-foreground">
                              {member.email}
                            </span>
                          </div>
                        </TableCell>
                        <TableCell>{member.student_id}</TableCell>
                        <TableCell>{member.course || "-"}</TableCell>
                        <TableCell>{member.contact_number || "-"}</TableCell>
                        <TableCell>
                          <Badge
                            variant="outline"
                            className={statusBadgeClass(
                              member.registration_status,
                            )}
                          >
                            {titleCaseStatus(member.registration_status)}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          <Badge
                            variant="outline"
                            className={statusBadgeClass(
                              member.latest_payment_status,
                            )}
                          >
                            {titleCaseStatus(member.latest_payment_status)}
                          </Badge>
                        </TableCell>
                        <TableCell>
                          {formatDateTime(member.created_at)}
                        </TableCell>
                      </TableRow>
                    ))}
                  </>
                )}
              </TableBody>
              {!isPending &&
                !isFetching &&
                !errorMessage &&
                totalMembers > 0 && (
                  <TableFooter>
                    <TableRow>
                      <TableCell colSpan={7} className="p-0">
                        <div className="flex flex-col gap-2 p-4 md:flex-row md:items-center md:justify-between">
                          <div className="flex items-center gap-3">
                            <p className="text-muted-foreground">
                              Showing {visibleStart}-{visibleEnd} of{" "}
                              {totalMembers} members
                            </p>
                          </div>
                          <div className="flex flex-wrap items-center gap-2">
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() =>
                                setCurrentPage((prev) => Math.max(1, prev - 1))
                              }
                              disabled={currentPage === 1}
                            >
                              Previous
                            </Button>
                            <span className="px-1 text-muted-foreground">
                              Page {currentPage} of {totalPages}
                            </span>
                            <Button
                              variant="outline"
                              size="sm"
                              onClick={() =>
                                setCurrentPage((prev) =>
                                  Math.min(totalPages, prev + 1),
                                )
                              }
                              disabled={safeCurrentPage === totalPages}
                            >
                              Next
                            </Button>
                          </div>
                          <Select
                            value={String(pageSize)}
                            onValueChange={(value) => {
                              setCurrentPage(1);
                              setPageSize(Number(value));
                            }}
                          >
                            <SelectTrigger className="h-8 w-28">
                              <SelectValue placeholder="Rows" />
                            </SelectTrigger>
                            <SelectContent>
                              <SelectItem value="10">10 / page</SelectItem>
                              <SelectItem value="25">25 / page</SelectItem>
                              <SelectItem value="50">50 / page</SelectItem>
                            </SelectContent>
                          </Select>
                        </div>
                      </TableCell>
                    </TableRow>
                  </TableFooter>
                )}
            </Table>
          </CardContent>
        </Card>
      </section>
    </DashboardLayout>
  );
}
