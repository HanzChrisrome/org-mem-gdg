import { BrowserRouter, Route, Routes } from "react-router-dom";

import LoginPage from "../pages/auth/LoginPage";
import ExecutiveDashboard from "../pages/dashboard/ExecutiveDashboard";
import MembersPage from "../pages/members/MembersPage";
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

        <Route path="/members" element={<MembersPage />} />
      </Routes>
    </BrowserRouter>
  );
}
