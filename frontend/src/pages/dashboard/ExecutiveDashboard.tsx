import DashboardLayout from "@/components/layout/DashboardLayout";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import {
  ArrowRightIcon,
  CheckCircle2Icon,
  CircleAlertIcon,
  ClipboardListIcon,
  ShieldCheckIcon,
  User,
  UserPlusIcon,
  UsersIcon,
  UsersRound,
  Wallet,
} from "lucide-react";

const overviewCards = [
  {
    title: "Total members",
    value: "248",
    detail: "32 new self-registrations this month",
    status: "+14 this week",
    icon: UsersIcon,
  },
  {
    title: "Pending payment reviews",
    value: "18",
    detail: "7 submissions need review within 24 hours",
    status: "Priority queue",
    icon: CircleAlertIcon,
  },
  {
    title: "Approved memberships",
    value: "193",
    detail: "78% of registered members are fully cleared",
    status: "On track",
    icon: CheckCircle2Icon,
  },
  {
    title: "Audit log coverage",
    value: "96%",
    detail: "Recent actions recorded across member and payment flows",
    status: "Healthy",
    icon: ShieldCheckIcon,
  },
] as const;

const pipeline = [
  {
    label: "New registrations",
    count: 32,
    width: "w-[82%]",
    tone: "bg-sky-500/80",
    note: "Awaiting initial profile validation",
  },
  {
    label: "Payment proofs pending",
    count: 18,
    width: "w-[54%]",
    tone: "bg-amber-500/80",
    note: "Executive review required",
  },
  {
    label: "Approved members",
    count: 193,
    width: "w-full",
    tone: "bg-emerald-500/80",
    note: "Eligible for events and reports",
  },
  {
    label: "Resubmissions",
    count: 5,
    width: "w-[26%]",
    tone: "bg-rose-500/80",
    note: "Need corrected proof or updated details",
  },
] as const;

const recentActivities = [
  {
    title: "Payment Queue Updated",
    description: "Proof #2719 was marked for executive review",
    time: "11:20 AM",
    icon: ClipboardListIcon,
    iconStyle: "bg-blue-100 text-blue-600",
  },
  {
    title: "New Member Registered",
    description: "Juan Dela Cruz completed self-registration",
    time: "11:15 AM",
    icon: UserPlusIcon,
    iconStyle: "bg-emerald-100 text-emerald-600",
  },
  {
    title: "Resubmission Received",
    description: "Proof #2322 was replaced after rejection feedback",
    time: "11:00 AM",
    icon: CircleAlertIcon,
    iconStyle: "bg-amber-100 text-amber-600",
  },
  {
    title: "Membership Approved",
    description: "Alyssa Dela Cruz is now marked as active",
    time: "10:45 AM",
    icon: CheckCircle2Icon,
    iconStyle: "bg-emerald-100 text-emerald-600",
  },
  {
    title: "Audit Export Generated",
    description: "Finance and membership logs exported for review",
    time: "10:30 AM",
    icon: ShieldCheckIcon,
    iconStyle: "bg-cyan-100 text-cyan-600",
  },
] as const;

export default function ExecutiveDashboard() {
  return (
    <DashboardLayout>
      <section className="space-y-4">
        <Card className="border-none bg-linear-to-br from-primary to-primary/80 text-primary-foreground shadow-sm">
          <CardHeader>
            <div className="flex flex-wrap items-center gap-2">
              <Badge className="bg-primary-foreground/15 text-primary-foreground hover:bg-primary-foreground/15">
                GDG on Campus Operations
              </Badge>
              <Badge className="bg-primary-foreground/10 text-primary-foreground hover:bg-primary-foreground/10">
                Executive dashboard
              </Badge>
            </div>
            <CardTitle className="text-3xl font-semibold tracking-tight text-primary-foreground pt-5">
              Welcome back, keep member onboarding, payment approvals, and
              reporting in one view.
            </CardTitle>
            <CardDescription className="text-primary-foreground/80">
              This workspace is tailored for organization executives who need to
              review registrations, validate payment proofs, monitor membership
              health, and maintain clean audit trails.
            </CardDescription>
          </CardHeader>
          <CardFooter className="mt-2 flex flex-col items-start justify-between gap-3 border-muted-foreground bg-transparent sm:flex-row sm:items-center">
            <div className="text-sm text-primary-foreground/80">
              18 payment submissions are waiting for action and 5 members need
              resubmission follow-up.
            </div>
            <div className="flex flex-wrap gap-2">
              <Button
                variant="secondary"
                size="lg"
                className="bg-primary-foreground text-primary hover:bg-primary-foreground/90"
              >
                <Wallet /> Review payments
              </Button>
              <Button
                variant="outline"
                size="lg"
                className="border-primary-foreground/25 bg-transparent text-primary-foreground hover:bg-primary-foreground/10 hover:text-primary-foreground"
              >
                <UsersRound /> Add new member
              </Button>
            </div>
          </CardFooter>
        </Card>

        <div className="grid grid-cols-1 gap-4 @xl/main:grid-cols-2 @5xl/main:grid-cols-4">
          {overviewCards.map(({ title, value, detail, status, icon: Icon }) => (
            <Card
              key={title}
              className="h-full bg-linear-to-b from-muted/30 to-card shadow-xs"
            >
              <CardHeader>
                <div className="flex items-start justify-between gap-3">
                  <div>
                    <CardDescription>{title}</CardDescription>
                    <CardTitle className="mt-2 text-3xl font-semibold tabular-nums">
                      {value}
                    </CardTitle>
                  </div>
                  <div className="rounded-xl bg-primary/10 p-2 text-primary">
                    <Icon className="size-5" />
                  </div>
                </div>
              </CardHeader>
              <CardFooter className="flex items-center justify-between gap-3 text-sm">
                <span className="text-muted-foreground">{detail}</span>
                <Badge variant="outline">{status}</Badge>
              </CardFooter>
            </Card>
          ))}
        </div>

        <div className="grid grid-cols-1 gap-4 @4xl/main:grid-cols-[1.4fr_1fr]">
          <Card>
            <CardHeader>
              <div className="flex items-start justify-between gap-3">
                <div>
                  <CardTitle>Membership pipeline</CardTitle>
                  <CardDescription>
                    Track where members are in the registration and payment
                    approval flow.
                  </CardDescription>
                </div>
                <div className="rounded-xl border p-2 text-primary">
                  <User className="size-5" />
                </div>
              </div>
            </CardHeader>
            <CardContent className="space-y-4 border-t pt-4">
              {pipeline.map((item) => (
                <div key={item.label} className="space-y-4">
                  <div className="flex items-center justify-between gap-3">
                    <div>
                      <p className="font-medium">{item.label}</p>
                      <p className="text-sm text-muted-foreground">
                        {item.note}
                      </p>
                    </div>
                    <Badge variant="outline">{item.count}</Badge>
                  </div>
                  <div className="h-2.5 rounded-full bg-muted">
                    <div
                      className={`h-2.5 rounded-full ${item.width} ${item.tone}`}
                    />
                  </div>
                </div>
              ))}
            </CardContent>
            <CardFooter className="justify-between gap-3 text-sm">
              <span className="text-muted-foreground">
                Last sync: 5 minutes ago
              </span>
              <Button variant="ghost" size="sm">
                View member request
                <ArrowRightIcon className="size-4" />
              </Button>
            </CardFooter>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle>Recent activity log</CardTitle>
              <CardDescription>
                Latest registration, payment, and audit actions from the
                executive workflow.
              </CardDescription>
            </CardHeader>
            <CardContent>
              <div className="border rounded-2xl h-[330px] overflow-auto">
                {recentActivities.map((item) => (
                  <div
                    key={item.title}
                    className="flex items-start justify-between gap-3 px-4 py-4 not-last:border-b"
                  >
                    <div className="flex min-w-0 items-start gap-3">
                      <div className={`rounded-xl p-2 ${item.iconStyle}`}>
                        <item.icon className="size-4" />
                      </div>
                      <div>
                        <p className="font-medium">{item.title}</p>
                        <p className="mt-1 text-sm text-muted-foreground">
                          {item.description}
                        </p>
                      </div>
                    </div>
                    <span className="shrink-0 text-xs text-muted-foreground">
                      {item.time}
                    </span>
                  </div>
                ))}
              </div>
            </CardContent>
            <CardFooter className="justify-between gap-3 text-sm h-full">
              <span className="text-muted-foreground">
                Auto-refresh every 60 seconds
              </span>
              <Button variant="ghost" size="sm">
                View full audit trail
                <ArrowRightIcon className="size-4" />
              </Button>
            </CardFooter>
          </Card>
        </div>
      </section>
    </DashboardLayout>
  );
}
