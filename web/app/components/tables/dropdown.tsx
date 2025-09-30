import { useState } from "react";
import { Link } from "react-router";

import { MoreHorizontal } from "lucide-react";

import { Button } from "~/components/ui/button";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "~/components/ui/dropdown-menu";

interface TableRowDropdownProps {
  editUrl: string;
  deleteUrl: string;
}

export default function TableRowDropdown({
  editUrl,
  deleteUrl,
}: TableRowDropdownProps) {
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" className="h-8 w-8 p-0 cursor-pointer">
          <span className="sr-only">Open menu</span>
          <MoreHorizontal />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end">
        <DropdownMenuItem asChild className="cursor-pointer">
          <Link to={editUrl}>Edit</Link>
        </DropdownMenuItem>
        <DropdownMenuItem asChild className="text-destructive cursor-pointer">
          <Link to={deleteUrl}>Delete</Link>
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
