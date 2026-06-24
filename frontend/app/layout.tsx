import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Maxy Sports — Calorie Tracker",
  description: "Track your calories and macros with Telegram auth",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en" suppressHydrationWarning>
      <body>{children}</body>
    </html>
  );
}
