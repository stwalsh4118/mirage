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
import { toast } from "sonner";
import {
  CheckCircle2,
  XCircle,
  AlertCircle,
  Loader2,
  Eye,
  EyeOff,
  ExternalLink,
  RotateCw,
  Trash2,
} from "lucide-react";
import {
  useRailwayTokenStatus,
  useStoreRailwayToken,
  useValidateRailwayToken,
  useDeleteRailwayToken,
  useRotateRailwayToken,
} from "@/hooks/useRailwayToken";

export function RailwayTokenTab() {
  const [addModalOpen, setAddModalOpen] = useState(false);
  const [rotateModalOpen, setRotateModalOpen] = useState(false);
  const [token, setToken] = useState("");
  const [newToken, setNewToken] = useState("");
  const [showToken, setShowToken] = useState(false);
  const [showNewToken, setShowNewToken] = useState(false);

  // API hooks
  const { data: status, isLoading: statusLoading } = useRailwayTokenStatus();
  const storeToken = useStoreRailwayToken();
  const validateToken = useValidateRailwayToken();
  const deleteToken = useDeleteRailwayToken();
  const rotateToken = useRotateRailwayToken();

  // Handle storing/updating token
  const handleStoreToken = async () => {
    if (!token.trim()) {
      toast.error("Token cannot be empty");
      return;
    }

    try {
      const response = await storeToken.mutateAsync({ token: token.trim() });
      toast.success(response.message || "Railway token saved successfully");
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
        toast.success(response.message || "Railway token is valid");
      } else {
        toast.error(response.message || "Railway token is invalid");
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
      toast.success(response.message || "Railway token deleted successfully");
    } catch (error) {
      const message = error instanceof Error ? error.message : "Failed to delete token";
      toast.error(message);
    }
  };

  // Handle rotating token
  const handleRotateToken = async () => {
    if (!newToken.trim()) {
      toast.error("New token cannot be empty");
      return;
    }

    try {
      const response = await rotateToken.mutateAsync({ new_token: newToken.trim() });
      toast.success(response.message || "Railway token rotated successfully");
      setNewToken("");
      setRotateModalOpen(false);
    } catch (error) {
      const message = error instanceof Error ? error.message : "Failed to rotate token";
      toast.error(message);
    }
  };

  // Get status display
  const getStatusDisplay = () => {
    if (statusLoading) {
      return {
        icon: <Loader2 className="h-5 w-5 animate-spin text-muted-foreground" />,
        text: "Checking status...",
        color: "text-muted-foreground",
      };
    }

    if (!status?.configured) {
      return {
        icon: <XCircle className="h-5 w-5 text-destructive" />,
        text: "Not Configured",
        color: "text-destructive",
      };
    }

    return {
      icon: <CheckCircle2 className="h-5 w-5 text-green-500" />,
      text: "Configured",
      color: "text-green-500",
    };
  };

  const statusDisplay = getStatusDisplay();

  return (
    <div className="space-y-4">
      {/* Status Card */}
      <Card>
        <CardHeader>
          <div className="flex items-start justify-between">
            <div>
              <CardTitle>Railway API Token</CardTitle>
              <CardDescription>
                Configure your Railway API token to enable environment provisioning and management
              </CardDescription>
            </div>
            <div className="flex items-center gap-2">
              {statusDisplay.icon}
              <span className={`text-sm font-medium ${statusDisplay.color}`}>
                {statusDisplay.text}
              </span>
            </div>
          </div>
        </CardHeader>
        <CardContent className="space-y-4">
          {/* Instructions */}
          {!status?.configured && (
            <div className="rounded-lg border border-blue-500/20 bg-blue-500/10 p-4">
              <div className="flex items-start gap-3">
                <AlertCircle className="h-5 w-5 text-blue-500 mt-0.5" />
                <div className="space-y-2 flex-1">
                  <p className="text-sm font-medium">How to get your Railway API token:</p>
                  <ol className="text-sm text-muted-foreground space-y-1 list-decimal list-inside">
                    <li>Go to Railway Account Settings</li>
                    <li>Navigate to the "Tokens" section</li>
                    <li>Create a new token or copy an existing one</li>
                    <li>Paste it below to get started</li>
                  </ol>
                  <a
                    href="https://railway.app/account/tokens"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-1 text-sm text-blue-500 hover:underline"
                  >
                    Open Railway Token Settings
                    <ExternalLink className="h-3 w-3" />
                  </a>
                </div>
              </div>
            </div>
          )}

          {/* Last Validated */}
          {status?.configured && status.last_validated && (
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <span>Last validated:</span>
              <span className="font-medium">
                {new Date(status.last_validated).toLocaleString()}
              </span>
            </div>
          )}

          {/* Action Buttons */}
          <div className="flex flex-wrap gap-2">
            {/* Add/Update Token Dialog */}
            <Dialog open={addModalOpen} onOpenChange={setAddModalOpen}>
              <DialogTrigger asChild>
                <Button variant={status?.configured ? "outline" : "default"}>
                  {status?.configured ? "Update Token" : "Add Token"}
                </Button>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>
                    {status?.configured ? "Update Railway Token" : "Add Railway Token"}
                  </DialogTitle>
                  <DialogDescription>
                    Enter your Railway API token. It will be validated before saving.
                  </DialogDescription>
                </DialogHeader>
                <div className="space-y-4 py-4">
                  <div className="space-y-2">
                    <Label htmlFor="token">Railway API Token</Label>
                    <div className="relative">
                      <Input
                        id="token"
                        type={showToken ? "text" : "password"}
                        value={token}
                        onChange={(e) => setToken(e.target.value)}
                        placeholder="Enter your Railway API token"
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
                </div>
                <DialogFooter>
                  <Button
                    variant="outline"
                    onClick={() => {
                      setAddModalOpen(false);
                      setToken("");
                    }}
                  >
                    Cancel
                  </Button>
                  <Button
                    onClick={handleStoreToken}
                    disabled={storeToken.isPending || !token.trim()}
                  >
                    {storeToken.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                    Save & Validate
                  </Button>
                </DialogFooter>
              </DialogContent>
            </Dialog>

            {/* Test Connection Button */}
            {status?.configured && (
              <Button
                variant="outline"
                onClick={handleValidateToken}
                disabled={validateToken.isPending}
              >
                {validateToken.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                Test Connection
              </Button>
            )}

            {/* Rotate Token Dialog */}
            {status?.configured && (
              <Dialog open={rotateModalOpen} onOpenChange={setRotateModalOpen}>
                <DialogTrigger asChild>
                  <Button variant="outline">
                    <RotateCw className="mr-2 h-4 w-4" />
                    Rotate Token
                  </Button>
                </DialogTrigger>
                <DialogContent>
                  <DialogHeader>
                    <DialogTitle>Rotate Railway Token</DialogTitle>
                    <DialogDescription>
                      Replace your current token with a new one. The old version will be preserved
                      for audit purposes.
                    </DialogDescription>
                  </DialogHeader>
                  <div className="space-y-4 py-4">
                    <div className="space-y-2">
                      <Label htmlFor="new-token">New Railway API Token</Label>
                      <div className="relative">
                        <Input
                          id="new-token"
                          type={showNewToken ? "text" : "password"}
                          value={newToken}
                          onChange={(e) => setNewToken(e.target.value)}
                          placeholder="Enter your new Railway API token"
                          className="pr-10"
                        />
                        <Button
                          type="button"
                          variant="ghost"
                          size="sm"
                          className="absolute right-0 top-0 h-full px-3 hover:bg-transparent"
                          onClick={() => setShowNewToken(!showNewToken)}
                        >
                          {showNewToken ? (
                            <EyeOff className="h-4 w-4" />
                          ) : (
                            <Eye className="h-4 w-4" />
                          )}
                        </Button>
                      </div>
                    </div>
                  </div>
                  <DialogFooter>
                    <Button
                      variant="outline"
                      onClick={() => {
                        setRotateModalOpen(false);
                        setNewToken("");
                      }}
                    >
                      Cancel
                    </Button>
                    <Button
                      onClick={handleRotateToken}
                      disabled={rotateToken.isPending || !newToken.trim()}
                    >
                      {rotateToken.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                      Rotate Token
                    </Button>
                  </DialogFooter>
                </DialogContent>
              </Dialog>
            )}

            {/* Delete Token Alert Dialog */}
            {status?.configured && (
              <AlertDialog>
                <AlertDialogTrigger asChild>
                  <Button variant="destructive" className="ml-auto">
                    <Trash2 className="mr-2 h-4 w-4" />
                    Remove Token
                  </Button>
                </AlertDialogTrigger>
                <AlertDialogContent>
                  <AlertDialogHeader>
                    <AlertDialogTitle>Remove Railway Token?</AlertDialogTitle>
                    <AlertDialogDescription>
                      This will remove your Railway API token from secure storage. You will need to
                      add a new token to continue managing Railway environments.
                      <br />
                      <br />
                      This action cannot be undone, but you can always add the token again later.
                    </AlertDialogDescription>
                  </AlertDialogHeader>
                  <AlertDialogFooter>
                    <AlertDialogCancel>Cancel</AlertDialogCancel>
                    <AlertDialogAction
                      onClick={handleDeleteToken}
                      disabled={deleteToken.isPending}
                      className="bg-destructive hover:bg-destructive/90"
                    >
                      {deleteToken.isPending && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
                      Remove Token
                    </AlertDialogAction>
                  </AlertDialogFooter>
                </AlertDialogContent>
              </AlertDialog>
            )}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}

