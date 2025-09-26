import { Link } from "react-router";

import { Button } from "~/components/ui/button";

export default function Page() {
  return (
    <div className="h-full flex items-center justify-center">
      <div className="h-full max-w-[90rem] flex-1 flex flex-col p-8 gap-4">
        <h1 className="text-2xl font-bold mb-4">Books</h1>
        <Button className="w-24" asChild>
          <Link to="/books/create">Create</Link>
        </Button>
      </div>
    </div>
  );
}
