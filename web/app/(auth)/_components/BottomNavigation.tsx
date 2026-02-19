"use client";

import { Icon } from "@iconify/react";
import Link from "next/link";
import { usePathname } from "next/navigation";

const navItems = [
  {
    label: "ライブラリ",
    href: "/",
    icon: "material-symbols:book",
  },
  {
    label: "貸出・返却",
    href: "/lending",
    icon: "material-symbols:swap-horiz",
  },
] as const;

export const BottomNavigation: React.FC = () => {
  const pathname = usePathname();

  return (
    <nav
      className={[
        "fixed",
        "bottom-0",
        "left-0",
        "right-0",
        "z-50",
        "border-t",
        "border-slate-200",
        "bg-white",
        "dark:border-slate-800",
        "dark:bg-slate-900",
        "md:hidden",
      ].join(" ")}
    >
      <div
        className={[
          "mx-auto",
          "flex",
          "max-w-4xl",
          "items-center",
          "justify-around",
        ].join(" ")}
      >
        {navItems.map((item) => {
          const isActive = pathname === item.href;
          return (
            <Link
              key={item.href}
              href={item.href}
              aria-current={isActive ? "page" : undefined}
              className={[
                "flex",
                "flex-1",
                "flex-col",
                "items-center",
                "gap-1",
                "py-2",
                "transition-colors",
                isActive
                  ? "text-sky-700 dark:text-sky-500"
                  : "text-slate-500 hover:text-slate-700 dark:text-slate-400 dark:hover:text-slate-300",
              ].join(" ")}
            >
              <Icon icon={item.icon} className="size-6" />
              <span className="text-xs font-medium">{item.label}</span>
            </Link>
          );
        })}
      </div>
    </nav>
  );
};
