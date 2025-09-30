"use client";

import Image from "next/image";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { CreateEnvironmentDialog } from "@/components/wizard/CreateEnvironmentDialog";
import { Input } from "@/components/ui/input";
import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { Separator } from "@/components/ui/separator";

export function DashboardHeader() {
  return (
    <div className="glass grain sticky top-0 z-40">
      <div className="max-w-screen-2xl mx-auto px-8 py-3 flex items-center gap-3">
        <div className="flex items-center gap-2 pr-2">
          <Link href="/dashboard" aria-label="Go to Dashboard" className="inline-flex items-center">
            <Image src="/mirage_logo.png" alt="Mirage" width={36} height={36} className="h-9 w-auto" />
          </Link>
        </div>
        <Separator orientation="vertical" className="mx-1 h-6" />
        <div className="flex-1">
          <Input placeholder="Search environments…  (⌘K)" className="h-9" />
        </div>
        <div className="flex items-center gap-2 pl-2">
          <CreateEnvironmentDialog trigger={<Button size="sm">New Environment</Button>} />
          <Avatar className="h-8 w-8">
            <AvatarImage alt="profile" />
            <AvatarFallback>ME</AvatarFallback>
          </Avatar>
        </div>
      </div>
    </div>
  );
}


