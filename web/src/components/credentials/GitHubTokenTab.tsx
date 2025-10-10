"use client";

import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "@/components/ui/alert-dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Badge } from "@/components/ui/badge";
import { toast } from "sonner";
import {
  CheckCircle2,
  XCircle,
  AlertCircle,
  Loader2,
  Eye,
  EyeOff,
  ExternalLink,
  Trash2,
  RefreshCw,
} from "lucide-react";
import {
  useGitHubTokenStatus,
  useStoreGitHubToken,
  useValidateGitHubToken,
  useDeleteGitHubToken,
} from "@/hooks/useGitHubToken";

export function GitHubTokenTab() {
  const [addModalOpen, setAddModalOpen] = useState(false);
  const [token, setToken] = useState("");
  const [showToken, setShowToken] = useState(false);

  // API hooks
  const { data: status, isLoading: statusLoading } = useGitHubTokenStatus();
  const storeToken = useStoreGitHubToken();
  const validateToken = useValidateGitHubToken();
  const deleteToken = useDeleteGitHubToken();

  // Handle storing/updating token
  const handleStoreToken = async () => {
    if (!token.trim()) {
      toast.error("Token cannot be empty");
      return;
    }

    try {
      const response = await storeToken.mutateAsync({ token: token.trim() });
      toast.success(response.message || "GitHub token saved successfully");
      setToken("");
      setAddModalOpen(false);
    } catch (error) {
      const message = error instanceof Error ? error.message : "Failed to save token";
      toast.error(message);
    }
  };

  // Handle validating token
  const handleValidateToken = async () => {
    try {
      const response = await validateToken.mutateAsync();
      if (response.valid) {
        toast.success(response.message || "GitHub token is valid");
      } else {
        toast.error(response.message || "GitHub token is invalid");
      }
    } catch (error) {
      const message = error instanceof Error ? error.message : "Failed to validate token";
      toast.error(message);
    }
  };

  // Handle deleting token
  const handleDeleteToken = async () => {
    try {
      const response = await deleteToken.mutateAsync();
      toast.success(response.message || "GitHub token deleted successfully");
    } catch (error) {
      const message = error instanceof Error ? error.message : "Failed to delete token";
      toast.error(message);
    }
  };

  // Loading state
  if (statusLoading) {
    return (
      <div className="flex items-center justify-center h-48">
        <Loader2 className="h-6 w-6 animate-spin text-muted-foreground" />
      </div>
    );
  }

  const isConfigured = status?.configured ?? false;

  return (
    <div className="space-y-4">
      {/* Status Card */}
      <Card>
        <CardHeader>
          <div className="flex items-center justify-between">
            <div className="space-y-1">
              <CardTitle>GitHub Personal Access Token</CardTitle>
              <CardDescription>
                Manage your GitHub token for accessing private repositories during Dockerfile discovery
              </CardDescription>
            </div>
            <div>
              {isConfigured ? (
                <CheckCircle2 className="h-5 w-5 text-green-500" />
              ) : (
                <XCircle className="h-5 w-5 text-muted-foreground" />
              )}
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {isConfigured ? (
            <>
              <div className="space-y-2">
                <div className="flex items-center gap-2">
                  <span className="text-sm text-muted-foreground">Status:</span>
                  <Badge variant="default" className="bg-green-500">
                    Configured
                  </Badge>
                </div>
                {status?.username && (
                  <div className="flex items-center gap-2">
                    <span className="text-sm text-muted-foreground">GitHub User:</span>
                    <span className="text-sm font-medium">{status.username}</span>
                  </div>
                )}
                {status?.scopes && status.scopes.length > 0 && (
                  <div className="space-y-1">
                    <span className="text-sm text-muted-foreground">Token Scopes:</span>
                    <div className="flex flex-wrap gap-1">
                      {status.scopes.map((scope) => (
                        <Badge key={scope} variant="secondary" className="text-xs">
                          {scope}
                        </Badge>
                      ))}
                    </div>
                  </div>
                )}
                {status?.last_validated && (
                  <div className="flex items-center gap-2">
                    <span className="text-sm text-muted-foreground">Last Validated:</span>
                    <span className="text-sm">
                      {new Date(status.last_validated).toLocaleString()}
                    </span>
                  </div>
                )}
              </div>

              {/* Action Buttons */}
              <div className="flex flex-wrap gap-2">
                <Dialog open={addModalOpen} onOpenChange={setAddModalOpen}>
                  <DialogTrigger asChild>
                    <Button variant="outline" size="sm">
                      Update Token
                    </Button>
                  </DialogTrigger>
                  <DialogContent>
                    <DialogHeader>
                      <DialogTitle>Update GitHub Token</DialogTitle>
                      <DialogDescription>
                        Replace your existing GitHub Personal Access Token. The old token will be
                        permanently removed.
                      </DialogDescription>
                    </DialogHeader>
                    <div className="space-y-4 py-4">
                      <div className="space-y-2">
                        <Label htmlFor="token">GitHub Personal Access Token</Label>
                        <div className="relative">
                          <Input
                            id="token"
                            type={showToken ? "text" : "password"}
                            value={token}
                            onChange={(e) => setToken(e.target.value)}
                            placeholder="ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
                            className="pr-10"
                          />
                          <Button
                            type="button"
                            variant="ghost"
                            size="sm"
                            className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
                            onClick={() => setShowToken(!showToken)}
                          >
                            {showToken ? (
                              <EyeOff className="h-4 w-4" />
                            ) : (
                              <Eye className="h-4 w-4" />
                            )}
                          </Button>
                        </div>
                        <p className="text-xs text-muted-foreground">
                          Your token is encrypted and stored securely in HashiCorp Vault
                        </p>
                      </div>
                      <div className="rounded-md bg-muted p-3 text-xs">
                        <p className="font-medium mb-1">Token Requirements:</p>
                        <ul className="list-disc list-inside space-y-1 text-muted-foreground">
                          <li>Must have <code className="text-xs">repo</code> scope for private repositories</li>
                          <li>Optionally <code className="text-xs">read:org</code> for organization repositories</li>
                        </ul>
                        <a
                          href="https://github.com/settings/tokens/new"
                          target="_blank"
                          rel="noopener noreferrer"
                          className="inline-flex items-center gap-1 mt-2 text-primary hover:underline"
                        >
                          Create a new token <ExternalLink className="h-3 w-3" />
                        </a>
                      </div>
                    </div>
                    <DialogFooter>
                      <Button
                        variant="outline"
                        onClick={() => {
                          setToken("");
                          setAddModalOpen(false);
                        }}
                        disabled={storeToken.isPending}
                      >
                        Cancel
                      </Button>
                      <Button
                        onClick={handleStoreToken}
                        disabled={!token.trim() || storeToken.isPending}
                      >
                        {storeToken.isPending && (
                          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        )}
                        Save Token
                      </Button>
                    </DialogFooter>
                  </DialogContent>
                </Dialog>

                <Button
                  variant="outline"
                  size="sm"
                  onClick={handleValidateToken}
                  disabled={validateToken.isPending}
                >
                  {validateToken.isPending ? (
                    <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                  ) : (
                    <RefreshCw className="mr-2 h-4 w-4" />
                  )}
                  Test Connection
                </Button>

                <AlertDialog>
                  <AlertDialogTrigger asChild>
                    <Button variant="outline" size="sm" className="text-destructive">
                      <Trash2 className="mr-2 h-4 w-4" />
                      Remove Token
                    </Button>
                  </AlertDialogTrigger>
                  <AlertDialogContent>
                    <AlertDialogHeader>
                      <AlertDialogTitle>Remove GitHub Token</AlertDialogTitle>
                      <AlertDialogDescription>
                        This will remove your GitHub token from Vault. You will need to add a new
                        token to access private repositories.
                      </AlertDialogDescription>
                    </AlertDialogHeader>
                    <AlertDialogFooter>
                      <AlertDialogCancel>Cancel</AlertDialogCancel>
                      <AlertDialogAction
                        onClick={handleDeleteToken}
                        disabled={deleteToken.isPending}
                        className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
                      >
                        {deleteToken.isPending && (
                          <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                        )}
                        Remove Token
                      </AlertDialogAction>
                    </AlertDialogFooter>
                  </AlertDialogContent>
                </AlertDialog>
              </div>
            </>
          ) : (
            <>
              <div className="flex items-start gap-3 rounded-lg border border-muted p-4">
                <AlertCircle className="h-5 w-5 text-muted-foreground flex-shrink-0 mt-0.5" />
                <div className="space-y-1">
                  <p className="text-sm font-medium">No GitHub Token Configured</p>
                  <p className="text-sm text-muted-foreground">
                    Add your GitHub Personal Access Token to enable Dockerfile discovery in
                    private repositories.
                  </p>
                </div>
              </div>

              {/* Add Token Dialog */}
              <Dialog open={addModalOpen} onOpenChange={setAddModalOpen}>
                <DialogTrigger asChild>
                  <Button>Add GitHub Token</Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Add GitHub Token</DialogTitle>
                    <DialogDescription>
                      Add your GitHub Personal Access Token to access private repositories during
                      Dockerfile discovery.
                    </DialogDescription>
                  </DialogHeader>
                  <div className="space-y-4 py-4">
                    <div className="space-y-2">
                      <Label htmlFor="token">GitHub Personal Access Token</Label>
                      <div className="relative">
                        <Input
                          id="token"
                          type={showToken ? "text" : "password"}
                          value={token}
                          onChange={(e) => setToken(e.target.value)}
                          placeholder="ghp_xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
                          className="pr-10"
                        />
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
                          onClick={() => setShowToken(!showToken)}
                        >
                          {showToken ? (
                            <EyeOff className="h-4 w-4" />
                          ) : (
                            <Eye className="h-4 w-4" />
                          )}
                        </Button>
                      </div>
                      <p className="text-xs text-muted-foreground">
                        Your token is encrypted and stored securely in HashiCorp Vault
                      </p>
                    </div>
                    <div className="rounded-md bg-muted p-3 text-xs">
                      <p className="font-medium mb-1">Token Requirements:</p>
                      <ul className="list-disc list-inside space-y-1 text-muted-foreground">
                        <li>Must have <code className="text-xs">repo</code> scope for private repositories</li>
                        <li>Optionally <code className="text-xs">read:org</code> for organization repositories</li>
                      </ul>
                      <a
                        href="https://github.com/settings/tokens/new"
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-flex items-center gap-1 mt-2 text-primary hover:underline"
                      >
                        Create a new token <ExternalLink className="h-3 w-3" />
                      </a>
                    </div>
                  </div>
                  <DialogFooter>
                    <Button
                      variant="outline"
                      onClick={() => {
                        setToken("");
                        setAddModalOpen(false);
                      }}
                      disabled={storeToken.isPending}
                    >
                      Cancel
                    </Button>
                    <Button
                      onClick={handleStoreToken}
                      disabled={!token.trim() || storeToken.isPending}
                    >
                      {storeToken.isPending && (
                        <Loader2 className="mr-2 h-4 w-4 animate-spin" />
                      )}
                      Save & Validate
                    </Button>
                  </DialogFooter>
                </DialogContent>
              </Dialog>
            </>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

