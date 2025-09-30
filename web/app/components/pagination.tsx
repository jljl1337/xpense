import { Link } from "react-router";

import {
  ArrowLeftToLineIcon,
  ArrowRightToLineIcon,
  ChevronLeftIcon,
  ChevronRightIcon,
} from "lucide-react";

import { Button } from "~/components/ui/button";

interface PaginationProps {
  page: number;
  totalPages: number;
  firstPageUrl: string;
  lastPageUrl: string;
  previousPageUrl: string;
  nextPageUrl: string;
}

export default function Pagination({
  page,
  totalPages,
  firstPageUrl,
  lastPageUrl,
  previousPageUrl,
  nextPageUrl,
}: PaginationProps) {
  return (
    <div className="ml-auto flex items-center gap-2 lg:ml-0">
      <Button
        variant="outline"
        className="hidden h-8 w-8 p-0 lg:flex"
        disabled={page === 1}
        asChild={!(page === 1)}
      >
        <Link to={firstPageUrl}>
          <span className="sr-only">Go to first page</span>
          <ArrowLeftToLineIcon />
        </Link>
      </Button>
      <Button
        variant="outline"
        className="size-8"
        size="icon"
        disabled={page === 1}
        asChild={!(page === 1)}
      >
        <Link to={previousPageUrl}>
          <span className="sr-only">Go to previous page</span>
          <ChevronLeftIcon />
        </Link>
      </Button>
      <div>
        <span className="text-sm text-muted-foreground">
          Page {page} of {totalPages}
        </span>
      </div>
      <Button
        variant="outline"
        className="size-8"
        size="icon"
        disabled={page === totalPages}
        asChild={!(page === totalPages)}
      >
        <Link to={nextPageUrl}>
          <span className="sr-only">Go to next page</span>
          <ChevronRightIcon />
        </Link>
      </Button>
      <Button
        variant="outline"
        className="hidden size-8 lg:flex"
        size="icon"
        disabled={page === totalPages}
        asChild={!(page === totalPages)}
      >
        <Link to={lastPageUrl}>
          <span className="sr-only">Go to last page</span>
          <ArrowRightToLineIcon />
        </Link>
      </Button>
    </div>
  );
}
