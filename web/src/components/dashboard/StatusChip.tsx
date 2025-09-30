import { cn } from "@/lib/utils";

type Status = "Running" | "Stopped" | "Creating" | "Destroying" | "Error" | "Unknown";

export function StatusChip({ status }: { status: Status }) {
  const { dot, text } = getStatusClasses(status);
  return (
    <span className={cn("inline-flex items-center gap-1.5 rounded-full px-2 py-0.5 text-xs font-medium", text)}>
      <span className={cn("h-1.5 w-1.5 rounded-full", dot)} />
      {status}
    </span>
  );
}

function getStatusClasses(status: Status) {
  switch (status) {
    case "Running":
      return { dot: "bg-green-500", text: "bg-green-500/10 text-green-700 dark:text-green-400" };
    case "Creating":
      return { dot: "bg-amber-500", text: "bg-amber-500/10 text-amber-700 dark:text-amber-400" };
    case "Destroying":
      return { dot: "bg-orange-500", text: "bg-orange-500/10 text-orange-700 dark:text-orange-400" };
    case "Error":
      return { dot: "bg-red-500", text: "bg-red-500/10 text-red-700 dark:text-red-400" };
    case "Stopped":
      return { dot: "bg-gray-400", text: "bg-gray-400/10 text-gray-700 dark:text-gray-300" };
    default:
      return { dot: "bg-teal-500", text: "bg-teal-500/10 text-teal-700 dark:text-teal-400" };
  }
}











