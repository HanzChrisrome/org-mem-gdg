import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

export default function LoginPage() {
  return (
    <div className="flex items-center justify-center h-screen bg-gray-100">
      <div className="w-[400px] bg-white p-6 rounded-lg shadow">
        <h1 className="text-2xl font-bold mb-6">Executive Login</h1>

        <div className="space-y-4">
          <Input placeholder="Email" />

          <Input type="password" placeholder="Password" />

          <Button className="w-full">Login</Button>
        </div>
      </div>
    </div>
  );
}
