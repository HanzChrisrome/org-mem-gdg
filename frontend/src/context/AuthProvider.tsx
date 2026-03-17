import { initAuth } from "@/api/auth";
import {
  createContext,
  type ReactNode,
  useContext,
  useEffect,
  useState,
} from "react";

interface AuthContextType {
  isLoggedIn: boolean;
  loading: boolean;
  setLoggedIn: (value: boolean) => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    initAuth()
      .then((valid) => setIsLoggedIn(valid))
      .finally(() => setLoading(false));
  }, []);

  return (
    <AuthContext.Provider
      value={{
        isLoggedIn,
        loading,
        setLoggedIn: setIsLoggedIn,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

// eslint-disable-next-line react-refresh/only-export-components
export function useAuth() {
  const context = useContext(AuthContext);
  if (!context) throw new Error("useAuth must be used within AuthProvider");
  return context;
}
