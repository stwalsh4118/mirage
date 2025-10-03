"use client";

import { useState, useCallback } from "react";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { Upload, FileText } from "lucide-react";

interface EnvImportDialogProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onImport: (variables: Array<{ key: string; value: string }>) => void;
  title?: string;
  description?: string;
}

export function EnvImportDialog({
  open,
  onOpenChange,
  onImport,
  title = "Import Environment Variables",
  description = "Paste your .env file contents or drag & drop a file",
}: EnvImportDialogProps) {
  const [content, setContent] = useState("");
  const [isDragging, setIsDragging] = useState(false);

  const parseEnvContent = (text: string): Array<{ key: string; value: string }> => {
    return text
      .split(/\r?\n/)
      .map((line) => line.trim())
      .filter((line) => line && !line.startsWith("#")) // Filter empty lines and comments
      .map((line) => {
        const eqIndex = line.indexOf("=");
        if (eqIndex === -1) {
          return { key: line, value: "" };
        }
        const key = line.slice(0, eqIndex).trim();
        let value = line.slice(eqIndex + 1).trim();
        
        // Remove surrounding quotes if present
        if ((value.startsWith('"') && value.endsWith('"')) || 
            (value.startsWith("'") && value.endsWith("'"))) {
          value = value.slice(1, -1);
        }
        
        return { key, value };
      })
      .filter((v) => v.key.length > 0);
  };

  const handleImport = () => {
    const variables = parseEnvContent(content);
    if (variables.length > 0) {
      onImport(variables);
      setContent("");
      onOpenChange(false);
    }
  };

  const handleFileRead = (file: File) => {
    // Limit file size to 1MB to prevent performance issues
    const MAX_FILE_SIZE = 1024 * 1024; // 1MB
    if (file.size > MAX_FILE_SIZE) {
      alert("File is too large. Maximum size is 1MB.");
      return;
    }

    const reader = new FileReader();
    reader.onload = (e) => {
      const text = e.target?.result as string;
      setContent(text);
    };
    reader.onerror = () => {
      alert("Failed to read file. Please try again.");
    };
    reader.readAsText(file);
  };

  const handleDrop = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);

    const files = Array.from(e.dataTransfer.files);
    const envFile = files.find((f) => f.name.endsWith(".env") || f.type === "text/plain");
    
    if (envFile) {
      handleFileRead(envFile);
    }
  }, []);
    
  const handleDragLeave = useCallback((e: React.DragEvent<HTMLDivElement>) => {
    e.preventDefault();
    // Only set isDragging to false if we're leaving the drop zone entirely
    if (!e.currentTarget.contains(e.relatedTarget as Node)) {
      setIsDragging(false);
    }
  }, []);

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  }, []);

  const handleFileSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (file) {
      handleFileRead(file);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="sm:max-w-[600px]">
        <DialogHeader>
          <DialogTitle>{title}</DialogTitle>
          <DialogDescription>{description}</DialogDescription>
        </DialogHeader>

        <div className="space-y-4 py-4">
          {/* Drag & Drop Zone */}
          <div
            onDrop={handleDrop}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            className={`
              relative border-2 border-dashed rounded-lg p-8 text-center transition-colors
              ${isDragging ? "border-primary bg-primary/5" : "border-border/60 bg-muted/20"}
            `}
          >
            <input
              type="file"
              id="env-file-input"
              accept=".env,text/plain"
              onChange={handleFileSelect}
              className="hidden"
            />
            <div className="flex flex-col items-center gap-2">
              <Upload className="h-8 w-8 text-muted-foreground" />
              <div className="text-sm">
                <label
                  htmlFor="env-file-input"
                  className="text-primary cursor-pointer hover:underline"
                >
                  Click to upload
                </label>
                <span className="text-muted-foreground"> or drag and drop</span>
              </div>
              <p className="text-xs text-muted-foreground">.env or text file</p>
            </div>
          </div>

          {/* Text Area */}
          <div className="space-y-2">
            <div className="flex items-center gap-2 text-sm text-muted-foreground">
              <FileText className="h-4 w-4" />
              <span>Or paste your .env file contents:</span>
            </div>
            <Textarea
              value={content}
              onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => setContent(e.target.value)}
              placeholder="NODE_ENV=production\nPORT=3000\nDATABASE_URL=postgres://..."
              className="font-mono text-xs min-h-[200px] bg-card"
            />
            {content && (
              <p className="text-xs text-muted-foreground">
                {parseEnvContent(content).length} variable(s) detected
              </p>
            )}
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            Cancel
          </Button>
          <Button onClick={handleImport} disabled={!content.trim()}>
            Import Variables
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}

