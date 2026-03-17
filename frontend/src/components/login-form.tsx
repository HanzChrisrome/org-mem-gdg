import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Field, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { cn } from "@/lib/utils";

import gdgoc_logo from "@/assets/gdgoc-logo.png";
import login_bg from "@/assets/login-bg.png";

export function LoginForm({
  className,
  ...props
}: React.ComponentProps<"div">) {
  return (
    <div className={cn("flex flex-col gap-6", className)} {...props}>
      <Card className="overflow-hidden p-0">
        <CardContent className="grid p-0 md:min-h-130 md:grid-cols-2">
          <form className="flex h-full items-center p-6 md:p-8">
            <div className="mx-auto w-full max-w-sm">
              <FieldGroup>
                <div className="flex flex-col items-center gap-1 text-center">
                  <img
                    src={gdgoc_logo}
                    alt="GDGOC logo"
                    className="mb-6 h-14 w-auto"
                  />
                  <h1 className="text-2xl font-bold">Welcome back</h1>
                  <p className="text-balance text-muted-foreground">
                    Login to your GDGOC - OrgMan Account
                  </p>
                </div>
                <Field>
                  <FieldLabel htmlFor="email">Email</FieldLabel>
                  <Input
                    id="email"
                    type="email"
                    placeholder="m@example.com"
                    required
                  />
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
                  <Input id="password" type="password" required />
                </Field>
                <Field>
                  <Button size="lg" type="submit">
                    Login
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
