"use client";

import { useWizardStore } from "@/store/wizard";
import { Label } from "@/components/ui/label";
import { Input } from "@/components/ui/input";
import { Alert, AlertDescription, AlertTitle } from "@/components/ui/alert";

export function StepSource() {
  const { repositoryUrl, repositoryBranch, setField } = useWizardStore();
  const repoProvided = repositoryUrl.trim().length > 0;
  const branchMissing = repoProvided && repositoryBranch.trim().length === 0;

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="space-y-2">
          <Label htmlFor="repo">Repository URL</Label>
          <Input
            id="repo"
            className="bg-card"
            placeholder="github.com/org/repo"
            value={repositoryUrl}
            onChange={(e) => setField("repositoryUrl", e.target.value)}
            spellCheck={false}
          />
        </div>
        <div className="space-y-2">
          <Label htmlFor="branch">Branch</Label>
          <Input
            id="branch"
            className="bg-card"
            placeholder="main"
            value={repositoryBranch}
            onChange={(e) => setField("repositoryBranch", e.target.value)}
            spellCheck={false}
          />
        </div>
      </div>

      {!repoProvided && (
        <Alert className="bg-muted/30 border-border/60">
          <AlertTitle>Optional source</AlertTitle>
          <AlertDescription>
            You can leave the repository blank if your template provisions placeholder services. Add it now to create initial services from a repo and branch.
          </AlertDescription>
        </Alert>
      )}

      {branchMissing && (
        <Alert className="bg-destructive/10 text-destructive-foreground border-destructive/40">
          <AlertTitle>Branch required</AlertTitle>
          <AlertDescription>
            Please enter a branch when a repository is provided.
          </AlertDescription>
        </Alert>
      )}
    </div>
  );
}


