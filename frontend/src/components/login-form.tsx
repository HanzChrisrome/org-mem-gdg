import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { cn } from "@/lib/utils";

import { login } from "@/api/auth";
import gdgoc_logo from "@/assets/gdgoc-logo.png";
import login_bg from "@/assets/login-bg.png";
import { useAuth } from "@/context/AuthProvider";
import { type LoginFormData, loginSchema } from "@/types/auth";
import { zodResolver } from "@hookform/resolvers/zod";
import { Loader2, Lock, MailIcon } from "lucide-react";
import { useForm } from "react-hook-form";
import { useNavigate } from "react-router-dom";
import { toast } from "sonner";
import { InputGroup, InputGroupAddon, InputGroupInput } from "./ui/input-group";

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginFormData>({
    resolver: zodResolver(loginSchema),
  });
  const { setLoggedIn } = useAuth();
  const navigate = useNavigate();

  const onSubmit = async (data: LoginFormData) => {
    try {
      await login({
        identifier: data.identifier,
        password: data.password,
      });

      setLoggedIn(true);
      toast.success("Login successful!");
      navigate("/dashboard");
    } catch (error) {
      const message =
        error instanceof Error
          ? error.message
          : "Login failed. Check credentials.";
      toast.error(message);
      console.error("Login error:", error);
    }
  };

  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card className="overflow-hidden p-0">
        <CardContent className="grid p-0 md:min-h-130 md:grid-cols-2">
          <form
            onSubmit={handleSubmit(onSubmit)}
            className="flex h-full items-center p-6 md:p-8"
          >
            <div className="mx-auto w-full max-w-sm">
              <FieldGroup>
                <div className="flex flex-col items-center text-center">
                  <img
                    src={gdgoc_logo}
                    alt="GDGOC logo"
                    className="mb-6 h-16 w-auto"
                  />
                  <h1 className="text-2xl font-bold">Welcome back</h1>
                  <p className="text-balance text-muted-foreground">
                    Login to your GDGOC - OrgMan Account
                  </p>
                </div>
                <Field>
                  <FieldLabel htmlFor="identifier">Email</FieldLabel>

                  <InputGroup>
                    <InputGroupInput
                      type="email"
                      placeholder="Enter your email"
                      {...register("identifier")}
                    />
                    <InputGroupAddon>
                      <MailIcon />
                    </InputGroupAddon>
                  </InputGroup>
                  {errors.identifier && (
                    <p className="text-sm text-red-500 mt-1">
                      {errors.identifier.message}
                    </p>
                  )}
                </Field>
                <Field>
                  <div className="flex items-center">
                    <FieldLabel htmlFor="password">Password</FieldLabel>
                    <a
                      href="#"
                      className="ml-auto text-sm underline-offset-2 hover:underline"
                    >
                      Forgot your password?
                    </a>
                  </div>
                  <InputGroup>
                    <InputGroupInput
                      id="password"
                      type="password"
                      placeholder="Enter your password"
                      {...register("password")}
                    />
                    <InputGroupAddon>
                      <Lock />
                    </InputGroupAddon>
                  </InputGroup>
                  {errors.password && (
                    <p className="text-sm text-red-500">
                      {errors.password.message}
                    </p>
                  )}
                </Field>
                <Field>
                  <Button
                    size="lg"
                    type="submit"
                    disabled={isSubmitting}
                    className="w-full"
                  >
                    {isSubmitting ? (
                      <>
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        Logging in...
                      </>
                    ) : (
                      "Login"
                    )}
                  </Button>
                </Field>
              </FieldGroup>
            </div>
          </form>
          <div className="hidden md:flex relative h-80 md:h-full">
            <img
              src={login_bg}
              alt="Login image"
              className="absolute inset-0 h-full w-full object-cover"
            />
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
