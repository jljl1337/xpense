export default function Error() {
  return (
    <div className="flex h-screen w-screen flex-col items-center justify-center bg-background">
      <h1 className="text-6xl font-bold">Oops!</h1>
      <p className="mt-4 text-xl">
        Something went wrong, please try again later
      </p>
    </div>
  );
}
