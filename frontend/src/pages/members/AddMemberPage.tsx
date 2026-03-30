import { createMember, updateMember } from "@/api/members";
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
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Textarea } from "@/components/ui/textarea";
import type { RegistrationStatus } from "@/types/member";
import { useMutation, useQueryClient } from "@tanstack/react-query";
import { type FormEvent, useState } from "react";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";

interface MemberFormState {
  name: string;
  email: string;
  student_id: string;
  password: string;
  course: string;
  contact_number: string;
  registration_status: RegistrationStatus;
  payment_status: string;
  payment_proof: string;
  notes: string;
}

const initialFormState: MemberFormState = {
  name: "",
  email: "",
  student_id: "",
  password: "",
  course: "",
  contact_number: "",
  registration_status: "pending",
  payment_status: "pending",
  payment_proof: "",
  notes: "",
};

function extractErrorMessage(error: unknown) {
  const apiError = error as {
    response?: {
      data?: { error?: string; message?: string };
    };
  };

  return (
    apiError.response?.data?.error ||
    apiError.response?.data?.message ||
    "Failed to create member. Please try again."
  );
}

export default function AddMemberPage() {
  const queryClient = useQueryClient();
  const navigate = useNavigate();
  const [form, setForm] = useState<MemberFormState>(initialFormState);

  const createMemberMutation = useMutation({
    mutationFn: async (payload: MemberFormState) => {
      const createdMember = await createMember({
        name: payload.name.trim(),
        email: payload.email.trim(),
        student_id: payload.student_id.trim(),
        password: payload.password,
        source_dashboard: "members",
      });

      const shouldUpdateExtraFields =
        Boolean(payload.course.trim()) ||
        Boolean(payload.contact_number.trim()) ||
        payload.registration_status !== "pending";

      if (shouldUpdateExtraFields) {
        await updateMember(createdMember.member_id, {
          course: payload.course.trim() || undefined,
          contact_number: payload.contact_number.trim() || undefined,
          registration_status: payload.registration_status,
        });
      }

      return createdMember;
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ["members"] });
      toast.success("Member created successfully.");
      navigate("/members");
    },
    onError: (error) => {
      toast.error(extractErrorMessage(error));
    },
  });

  const isSubmitting = createMemberMutation.isPending;

  const setField = <K extends keyof MemberFormState>(
    field: K,
    value: MemberFormState[K],
  ) => {
    setForm((prev) => ({ ...prev, [field]: value }));
  };

  async function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();

    if (!form.name || !form.email || !form.student_id || !form.password) {
      toast.error("Name, email, student ID, and password are required.");
      return;
    }

    await createMemberMutation.mutateAsync(form);
  }

  return (
    <DashboardLayout>
      <section className="flex flex-col gap-6 py-4 md:py-6">
        <div className="px-4 lg:px-6">
          <Card>
            <CardHeader>
              <div className="flex flex-wrap items-center gap-2">
                <CardTitle>Add New Member</CardTitle>
                <Badge
                  variant="outline"
                  className="bg-sky-100 text-sky-700 border-sky-200 hover:bg-sky-100"
                >
                  Connected form
                </Badge>
              </div>
              <CardDescription>
                Create a member using backend endpoints. Course, contact number,
                and registration status are saved after member creation.
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-6">
              <form className="space-y-6" onSubmit={handleSubmit}>
                <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
                  <div className="space-y-2">
                    <Label htmlFor="full-name">Full Name</Label>
                    <Input
                      id="full-name"
                      placeholder="Juan Dela Cruz"
                      value={form.name}
                      onChange={(event) => setField("name", event.target.value)}
                      required
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="email">Email Address</Label>
                    <Input
                      id="email"
                      type="email"
                      placeholder="juan.delacruz@school.edu"
                      value={form.email}
                      onChange={(event) =>
                        setField("email", event.target.value)
                      }
                      required
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="student-id">Student ID</Label>
                    <Input
                      id="student-id"
                      placeholder="2026-00012"
                      value={form.student_id}
                      onChange={(event) =>
                        setField("student_id", event.target.value)
                      }
                      required
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="password">Temporary Password</Label>
                    <Input
                      id="password"
                      type="password"
                      placeholder="At least 12 chars, upper/lower/digit/special"
                      value={form.password}
                      onChange={(event) =>
                        setField("password", event.target.value)
                      }
                      required
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="course">Course / Program</Label>
                    <Input
                      id="course"
                      placeholder="BS Computer Science"
                      value={form.course}
                      onChange={(event) =>
                        setField("course", event.target.value)
                      }
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="contact-number">Contact Number</Label>
                    <Input
                      id="contact-number"
                      placeholder="0917-000-0000"
                      value={form.contact_number}
                      onChange={(event) =>
                        setField("contact_number", event.target.value)
                      }
                    />
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="registration-status">
                      Registration Status
                    </Label>
                    <Select
                      value={form.registration_status}
                      onValueChange={(value) =>
                        setField(
                          "registration_status",
                          value as RegistrationStatus,
                        )
                      }
                    >
                      <SelectTrigger
                        id="registration-status"
                        className="w-full"
                      >
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="pending">Pending</SelectItem>
                        <SelectItem value="active">Active</SelectItem>
                        <SelectItem value="inactive">Inactive</SelectItem>
                        <SelectItem value="resubmitted">Resubmitted</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="payment-status">Payment Status</Label>
                    <Select
                      value={form.payment_status}
                      onValueChange={(value) =>
                        setField("payment_status", value)
                      }
                      disabled
                    >
                      <SelectTrigger id="payment-status" className="w-full">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="pending">Pending</SelectItem>
                        <SelectItem value="approved">Approved</SelectItem>
                        <SelectItem value="rejected">Rejected</SelectItem>
                        <SelectItem value="resubmitted">Resubmitted</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div className="space-y-2">
                    <Label htmlFor="payment-proof">Payment Proof URL</Label>
                    <Input
                      id="payment-proof"
                      placeholder="https://example.com/proof.png"
                      value={form.payment_proof}
                      onChange={(event) =>
                        setField("payment_proof", event.target.value)
                      }
                      disabled
                    />
                  </div>
                </div>

                <div className="space-y-2">
                  <Label htmlFor="notes">Executive Notes</Label>
                  <Textarea
                    id="notes"
                    placeholder="Optional notes for audit logging, review remarks, or verification comments."
                    value={form.notes}
                    onChange={(event) => setField("notes", event.target.value)}
                    disabled
                  />
                  <p className="text-xs text-muted-foreground">
                    Payment and notes fields are displayed for planning, but
                    they are not submitted because no create-payment endpoint is
                    available in the current Swagger member flow.
                  </p>
                </div>

                <div className="flex flex-wrap items-center gap-3">
                  <Button type="submit" disabled={isSubmitting}>
                    {isSubmitting ? "Creating member..." : "Create member"}
                  </Button>
                  <Button
                    type="button"
                    variant="outline"
                    onClick={() => setForm(initialFormState)}
                    disabled={isSubmitting}
                  >
                    Reset form
                  </Button>
                </div>
              </form>
            </CardContent>
          </Card>
        </div>
      </section>
    </DashboardLayout>
  );
}
