"use client";

import { useState } from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from "@/components/ui/card";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { MoreVertical, Trash2 } from "lucide-react";
import { RailwayProjectDetails } from "@/lib/api/railway";
import { useDeleteRailwayProject } from "@/hooks/useRailway";
import { toast } from "sonner";

export function ProjectCard({ project }: { project: RailwayProjectDetails }) {
  const servicesCount = project.services?.length ?? 0;
  const environmentsCount = project.environments?.length ?? 0;
  const [showDeleteDialog, setShowDeleteDialog] = useState(false);
  const deleteProject = useDeleteRailwayProject();

  const handleDelete = async () => {
    try {
      await deleteProject.mutateAsync(project.id);
      toast.success(`"${project.name}" has been permanently deleted from Railway.`);
      setShowDeleteDialog(false);
    } catch (error) {
      toast.error((error as Error).message || "Could not delete project. Please try again.");
    }
  };

  return (
    <>
      <Card className="glass grain transition-all duration-200 hover:translate-y-[-1px] hover:scale-[1.01]">
        <CardHeader className="flex flex-row items-start justify-between gap-2">
          <div className="space-y-1">
            <CardTitle className="text-base font-medium">{project.name}</CardTitle>
            <div className="text-[11px] text-muted-foreground">{project.id}</div>
          </div>
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" size="icon" className="h-7 w-7">
                <MoreVertical className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuItem
                className="text-destructive focus:text-destructive"
                onClick={() => setShowDeleteDialog(true)}
              >
                <Trash2 className="mr-2 h-4 w-4" />
                Delete project
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        </CardHeader>
        <CardContent className="space-y-3">
          <div className="flex items-center justify-between text-xs text-muted-foreground">
            <span>services</span>
            <span className="font-medium text-foreground/80">{servicesCount}</span>
          </div>
          <div className="flex items-center justify-between text-xs text-muted-foreground">
            <span>environments</span>
            <span className="font-medium text-foreground/80">{environmentsCount}</span>
          </div>
        </CardContent>
        <CardFooter className="flex items-center justify-between">
          <Button asChild variant="outline" size="sm" aria-label="Open project details" className="bg-transparent">
            <Link href={`/project/${project.id}`}>Open</Link>
          </Button>
        </CardFooter>
      </Card>

      <AlertDialog open={showDeleteDialog} onOpenChange={setShowDeleteDialog}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>Delete Railway Project?</AlertDialogTitle>
            <AlertDialogDescription className="space-y-2">
              <p>
                You are about to permanently delete <strong>{project.name}</strong>.
              </p>
              <p className="text-destructive font-medium">
                This will delete all {environmentsCount} environment(s), {servicesCount} service(s), and all associated data.
                This action cannot be undone.
              </p>
            </AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={deleteProject.isPending}>Cancel</AlertDialogCancel>
            <AlertDialogAction
              onClick={handleDelete}
              disabled={deleteProject.isPending}
              className="bg-destructive text-destructive-foreground hover:bg-destructive/90"
            >
              {deleteProject.isPending ? "Deleting..." : "Delete project"}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}


