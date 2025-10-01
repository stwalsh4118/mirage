"use client";

import { useState } from "react";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { AlertDialog, AlertDialogAction, AlertDialogCancel, AlertDialogContent, AlertDialogDescription, AlertDialogFooter, AlertDialogHeader, AlertDialogTitle } from "@/components/ui/alert-dialog";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuSeparator, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import Link from "next/link";
import type { RailwayEnvironmentWithServices } from "@/lib/api/railway";
import { useDeleteRailwayEnvironment } from "@/hooks/useRailway";
import { MoreVertical, Trash2 } from "lucide-react";
import { toast } from "sonner";
import { Pill } from "./Pill";

export function RailwayEnvironmentCard({ env, href }: { env: RailwayEnvironmentWithServices; href?: string }) {
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const deleteEnvironment = useDeleteRailwayEnvironment();
  const serviceCount = env.services?.length ?? 0;

  const handleDelete = async () => {
    try {
      await deleteEnvironment.mutateAsync(env.id);
      toast.success(`Environment "${env.name}" deleted successfully`);
      setShowDeleteDialog(false);
    } catch (error) {
      toast.error(`Failed to delete environment: ${error instanceof Error ? error.message : "Unknown error"}`);
    }
  };

  const content = (
    <Card className="glass grain transition-all duration-200 hover:translate-y-[-1px] hover:scale-[1.01]">
      <CardHeader className="flex flex-row items-start justify-between gap-2">
        <div className="space-y-1 flex-1">
          <CardTitle className="text-base font-medium">{env.name}</CardTitle>
          <div className="text-[11px] text-muted-foreground">{env.id}</div>
        </div>
        <div className="flex items-center gap-2">
          <Pill color={serviceCount > 0 ? "green" : "amber"}>
            {serviceCount} service{serviceCount === 1 ? "" : "s"}
          </Pill>
          <DropdownMenu>
            <DropdownMenuTrigger asChild onClick={(e) => e.preventDefault()}>
              <Button variant="ghost" size="icon" className="h-8 w-8" aria-label="More actions">
                <MoreVertical className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end" onClick={(e) => e.stopPropagation()}>
              {href && (
                <>
                  <DropdownMenuItem asChild>
                    <Link href={href}>View project</Link>
                  </DropdownMenuItem>
                  <DropdownMenuSeparator />
                </>
              )}
              <DropdownMenuItem
                className="text-destructive focus:text-destructive"
                onSelect={(e) => {
                  e.preventDefault();
                  setShowDeleteDialog(true);
                }}
              >
                <Trash2 className="h-4 w-4 mr-2" />
                Delete environment
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </div>
      </CardHeader>
      <CardContent className="space-y-3">
        <div className="text-xs text-muted-foreground">Railway environment</div>
        {serviceCount > 0 && (
          <div className="flex flex-wrap gap-2 text-[11px]">
            {env.services!.slice(0, 4).map((s) => (
              <Pill key={s.id} color="blue">{s.serviceName}</Pill>
            ))}
            {env.services!.length > 4 && (
              <Pill color="neutral">+{env.services!.length - 4} more</Pill>
            )}
          </div>
        )}
      </CardContent>
      <CardFooter className="flex items-center justify-between">
        <div className="text-[11px] text-muted-foreground">Railway</div>
        <div className="text-[11px] text-muted-foreground">{serviceCount} service{serviceCount === 1 ? "" : "s"}</div>
      </CardFooter>
    </Card>
  );

  return (
    <>
      {href ? (
        <Link href={href} className="block" aria-label={`Open project for ${env.name}`}>
          {content}
        </Link>
      ) : (
        content
      )}

      <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Railway environment?</AlertDialogTitle>
            <AlertDialogDescription>
              Are you sure you want to delete <strong>{env.name}</strong>? This action cannot be undone.
              {serviceCount > 0 && (
                <>
                  <br /><br />
                  This will destroy the Railway environment and all <strong>{serviceCount}</strong> associated service{serviceCount === 1 ? "" : "s"}.
                </>
              )}
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={deleteEnvironment.isPending}>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={(e) => {
                e.preventDefault();
                void handleDelete();
              }}
              disabled={deleteEnvironment.isPending}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {deleteEnvironment.isPending ? "Deleting..." : "Delete"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}


