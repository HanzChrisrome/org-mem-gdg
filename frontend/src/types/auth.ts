// src/types/auth.ts
import { z } from "zod";

export const loginSchema = z.object({
  identifier: z.email("Invalid email"),
  password: z.string().min(6, "Password must be at least 6 chars"),
});

export type LoginFormData = z.infer<typeof loginSchema>;
