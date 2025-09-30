import { cn } from "@/lib/utils";

type PillProps = {
  children: React.ReactNode;
  color?: "neutral" | "green" | "amber" | "red" | "blue" | "accent";
  className?: string;
};

export function Pill({ children, color = "neutral", className }: PillProps) {
  const colorClasses: Record<NonNullable<PillProps["color"]>, string> = {
    neutral: "bg-muted/60 text-foreground/80",
    green: "bg-green-500/10 text-green-700 dark:text-green-400",
    amber: "bg-amber-500/10 text-amber-700 dark:text-amber-400",
    red: "bg-red-500/10 text-red-700 dark:text-red-400",
    blue: "bg-blue-500/10 text-blue-700 dark:text-blue-400",
    accent: "bg-accent/15 text-accent",
  };
  return (
    <span
      className={cn(
        "inline-flex items-center rounded-full px-2.5 py-1 text-xs font-medium",
        colorClasses[color],
        className
      )}
    >
      {children}
    </span>
  );
}











