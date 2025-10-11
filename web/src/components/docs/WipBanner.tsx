import { AlertCircle } from "lucide-react";
import { Alert, AlertDescription } from "@/components/ui/alert";

export function WipBanner() {
  return (
    <Alert className="mb-6 border-amber-500/50 bg-amber-500/10">
      <AlertCircle className="h-4 w-4 text-amber-500" />
      <AlertDescription className="text-amber-600 dark:text-amber-400">
        <strong>Work in Progress:</strong> This documentation is currently being written and may be incomplete or subject to change.
      </AlertDescription>
    </Alert>
  );
}

