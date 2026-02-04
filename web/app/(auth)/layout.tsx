import type { Metadata } from "next";
import { BIZ_UDPGothic } from "next/font/google";
import "../globals.css";
import { AuthProvider } from "../_components/AuthProvider";
import { QueryProvider } from "../_components/QueryProvider";

const bizUDPGothic = BIZ_UDPGothic({
  variable: "--font-biz-udp-gothic",
  subsets: ["latin"],
  weight: ["400", "700"],
});

export const metadata: Metadata = {
  title: "Holocron",
  description: "Knowledge archive for the galaxy",
  metadataBase: process.env.NEXT_PUBLIC_BASE_URL
    ? new URL(process.env.NEXT_PUBLIC_BASE_URL)
    : undefined,
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      <body className={`${bizUDPGothic.variable} antialiased`}>
        <QueryProvider>
          <AuthProvider>{children}</AuthProvider>
        </QueryProvider>
      </body>
    </html>
  );
}
