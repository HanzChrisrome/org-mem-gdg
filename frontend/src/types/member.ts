export type RegistrationStatus =
  | "pending"
  | "active"
  | "inactive"
  | "resubmitted";

export type PaymentStatus = "pending" | "approved" | "rejected" | "resubmitted";

export interface MemberResponse {
  member_id: string;
  name: string;
  email: string;
  student_id: string;
  course?: string;
  contact_number?: string;
  registration_status: string;
  created_at: string;
  last_updated: string;
}

export interface MemberWithPaymentResponse extends MemberResponse {
  latest_payment_id?: string | null;
  latest_payment_status?: string | null;
  latest_submission_date?: string | null;
  latest_approval_date?: string | null;
  approver_name?: string | null;
}

export interface CreateMemberRequest {
  name: string;
  email: string;
  student_id: string;
  password: string;
  source_dashboard: "members";
}

export interface UpdateMemberRequest {
  name?: string;
  email?: string;
  student_id?: string;
  course?: string;
  contact_number?: string;
  registration_status?: string;
}
