import { redirect } from "next/navigation";

/**
 * Root `/` route — immediately sends the browser to `/login`.
 * The auth guard in each protected route will redirect to `/home` when a
 * valid token is already present.
 */
export default function RootPage() {
  redirect("/login");
}
