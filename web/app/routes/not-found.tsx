export default function NotFound() {
  // return a page takes up the full screen and says 404 - Not Found in the center
  return (
    <>
      <title>Page Not Found | Xpense</title>
      <div className="flex h-screen w-screen flex-col items-center justify-center bg-background">
        <h1 className="text-6xl font-bold">404</h1>
        <p className="mt-4 text-xl">Page Not Found</p>
      </div>
    </>
  );
}
