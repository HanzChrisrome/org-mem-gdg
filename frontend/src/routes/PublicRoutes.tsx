import { useAuth } from "@/context/AuthProvider";
import { Navigate } from "react-router-dom";

export default function PublicRoute({
  children,
}: {
  children: React.ReactNode;
}) {
  const { isLoggedIn, loading } = useAuth();

  if (loading) return <div>Loading...</div>;
  if (isLoggedIn) return <Navigate to="/dashboard" replace />;

  return children;
}
