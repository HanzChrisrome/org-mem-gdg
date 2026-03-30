import { BrowserRouter, Route, Routes } from "react-router-dom";

import LoginPage from "../pages/auth/LoginPage";
import ExecutiveDashboard from "../pages/dashboard/ExecutiveDashboard";
import AddMemberPage from "../pages/members/AddMemberPage";
import MembersPage from "../pages/members/MembersPage";
import PaymentApprovalPage from "../pages/payments/PaymentApprovalPage";
import ProtectedRoute from "./ProtectedRoutes";
import PublicRoute from "./PublicRoutes";

export default function AppRoutes() {
  return (
    <BrowserRouter>
      <Routes>
        <Route
          path="/login"
          element={
            <PublicRoute>
              <LoginPage />
            </PublicRoute>
          }
        />

        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <ExecutiveDashboard />
            </ProtectedRoute>
          }
        />

        <Route
          path="/members"
          element={
            <ProtectedRoute>
              <MembersPage />
            </ProtectedRoute>
          }
        />

        <Route
          path="/members/new"
          element={
            <ProtectedRoute>
              <AddMemberPage />
            </ProtectedRoute>
          }
        />

        <Route
          path="/payments/approvals"
          element={
            <ProtectedRoute>
              <PaymentApprovalPage />
            </ProtectedRoute>
          }
        />
      </Routes>
    </BrowserRouter>
  );
}
