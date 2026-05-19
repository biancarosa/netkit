import type { Metadata } from "next";
import "./globals.css";

const basePath = process.env.NEXT_PUBLIC_NETKIT_BASE_PATH || "";

export const metadata: Metadata = {
  title: "netkit - HTTP Requests Helper",
  description: "A powerful HTTP proxy and request capture tool with a modern web dashboard",
  icons: {
    icon: `${basePath}/favicon.ico`,
  },
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body className="antialiased">
        {children}
      </body>
    </html>
  );
}
