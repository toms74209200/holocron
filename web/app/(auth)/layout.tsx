import { BIZ_UDPGothic } from "next/font/google";
import "../globals.css";
import { AuthProvider } from "../_components/AuthProvider";
import { QueryProvider } from "../_components/QueryProvider";
import { BottomNavigation } from "./_components/BottomNavigation";
import { MobileHeader } from "./_components/MobileHeader";
import { SidebarNavigation } from "./_components/SidebarNavigation";

const bizUDPGothic = BIZ_UDPGothic({
  variable: "--font-biz-udp-gothic",
  subsets: ["latin"],
  weight: ["400", "700"],
});

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="ja">
      <body className={`${bizUDPGothic.variable} antialiased`}>
        <QueryProvider>
          <AuthProvider>
            <SidebarNavigation />
            <div className="md:ml-64">
              <MobileHeader />
              <div className="pb-16 md:pb-0">{children}</div>
            </div>
            <BottomNavigation />
          </AuthProvider>
        </QueryProvider>
      </body>
    </html>
  );
}
