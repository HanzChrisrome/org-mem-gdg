import type {
  CreateMemberRequest,
  MemberResponse,
  MemberWithPaymentResponse,
  UpdateMemberRequest,
} from "@/types/member";
import api from "./axios";

interface ListMembersParams {
  q?: string;
  status?: string;
}

export async function listMembers(params?: ListMembersParams) {
  const response = await api.get<MemberWithPaymentResponse[] | null>(
    "/members",
    {
      params,
    },
  );
  // Backend can return null when no rows match; normalize to [] for UI safety.
  return Array.isArray(response.data) ? response.data : [];
}

export async function createMember(payload: CreateMemberRequest) {
  const response = await api.post<MemberResponse>("/members", payload);
  return response.data;
}

export async function updateMember(
  memberId: string,
  payload: UpdateMemberRequest,
) {
  const response = await api.put<MemberResponse>(
    `/members/${memberId}`,
    payload,
  );
  return response.data;
}
