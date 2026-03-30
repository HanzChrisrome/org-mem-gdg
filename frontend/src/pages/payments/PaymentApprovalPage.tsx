import DashboardLayout from "@/components/layout/DashboardLayout";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
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
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";

const paymentSubmissions = [
  {
    payment_id: "PAY-2719",
    member_id: "MBR-1043",
    member_name: "Marco Sison",
    student_id: "2022-00981",
    payment_proof_image: "proof_2719.png",
    payment_status: "Pending",
    submission_date: "2026-03-23 07:50",
    approval_date: "-",
    approved_by: "-",
  },
  {
    payment_id: "PAY-2722",
    member_id: "MBR-1047",
    member_name: "Nica Velasco",
    student_id: "2025-00301",
    payment_proof_image: "proof_2722.jpg",
    payment_status: "Pending",
    submission_date: "2026-03-23 08:10",
    approval_date: "-",
    approved_by: "-",
  },
  {
    payment_id: "PAY-2322",
    member_id: "MBR-1044",
    member_name: "Janelle Reyes",
    student_id: "2024-00117",
    payment_proof_image: "proof_2322_v2.png",
    payment_status: "Resubmitted",
    submission_date: "2026-03-22 14:25",
    approval_date: "-",
    approved_by: "-",
  },
  {
    payment_id: "PAY-2693",
    member_id: "MBR-1042",
    member_name: "Alyssa Dela Cruz",
    student_id: "2023-01452",
    payment_proof_image: "proof_2693.png",
    payment_status: "Approved",
    submission_date: "2026-03-21 08:30",
    approval_date: "2026-03-21 10:05",
    approved_by: "EXEC-001",
  },
] as const;

function paymentBadgeClass(status: string) {
  if (status === "Approved") {
    return "bg-emerald-100 text-emerald-700 border-emerald-200 hover:bg-emerald-100";
  }

  if (status === "Pending") {
    return "bg-amber-100 text-amber-700 border-amber-200 hover:bg-amber-100";
  }

  return "bg-sky-100 text-sky-700 border-sky-200 hover:bg-sky-100";
}

export default function PaymentApprovalPage() {
  return (
    <DashboardLayout>
      <section className="flex flex-col gap-6 py-4 md:py-6">
        <div className="px-4 lg:px-6">
          <Card>
            <CardHeader className="pb-4">
              <CardTitle>Payment Approval</CardTitle>
              <CardDescription>
                Review payment submissions and verify proof before membership
                approval.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="flex flex-col gap-3 md:flex-row md:items-center md:justify-between">
                <div className="flex flex-1 flex-col gap-3 sm:flex-row">
                  <Input
                    placeholder="Search by member name, student ID, or payment ID"
                    className="sm:max-w-sm"
                  />
                  <Select defaultValue="all">
                    <SelectTrigger className="w-full sm:w-48">
                      <SelectValue placeholder="Filter by status" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="all">All statuses</SelectItem>
                      <SelectItem value="pending">Pending</SelectItem>
                      <SelectItem value="resubmitted">Resubmitted</SelectItem>
                      <SelectItem value="approved">Approved</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                <Button variant="outline">Export queue</Button>
              </div>

              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead>Payment ID</TableHead>
                    <TableHead>Member</TableHead>
                    <TableHead>Student ID</TableHead>
                    <TableHead>Proof</TableHead>
                    <TableHead>Status</TableHead>
                    <TableHead>Submitted</TableHead>
                    <TableHead>Approved By</TableHead>
                    <TableHead className="text-right">Actions</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {paymentSubmissions.map((payment) => (
                    <TableRow key={payment.payment_id}>
                      <TableCell className="font-medium">
                        {payment.payment_id}
                      </TableCell>
                      <TableCell>{payment.member_name}</TableCell>
                      <TableCell>{payment.student_id}</TableCell>
                      <TableCell>{payment.payment_proof_image}</TableCell>
                      <TableCell>
                        <Badge
                          variant="outline"
                          className={paymentBadgeClass(payment.payment_status)}
                        >
                          {payment.payment_status}
                        </Badge>
                      </TableCell>
                      <TableCell>{payment.submission_date}</TableCell>
                      <TableCell>{payment.approved_by}</TableCell>
                      <TableCell className="text-right">
                        <div className="flex justify-end gap-2">
                          <Button size="sm" variant="secondary">
                            Review
                          </Button>
                          <Button size="sm">Approve</Button>
                          <Button size="sm" variant="destructive">
                            Reject
                          </Button>
                        </div>
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </CardContent>
          </Card>
        </div>
      </section>
    </DashboardLayout>
  );
}
