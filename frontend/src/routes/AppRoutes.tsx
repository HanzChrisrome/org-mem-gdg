import { BrowserRouter, Route, Routes } from "react-router-dom";

import LoginPage from "../pages/auth/LoginPage";
import ExecutiveDashboard from "../pages/dashboard/ExecutiveDashboard";
import MembersPage from "../pages/members/MembersPage";

export default function AppRoutes() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />

        <Route path="/dashboard" element={<ExecutiveDashboard />} />

        <Route path="/members" element={<MembersPage />} />
      </Routes>
    </BrowserRouter>
  );
}
