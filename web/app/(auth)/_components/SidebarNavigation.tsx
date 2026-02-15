"use client";

import { Icon } from "@iconify/react";
import Image from "next/image";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { Fragment } from "react";
import holocronLogo from "./holocron-logo.svg";

const navItems = [
  {
    label: "ライブラリ",
    href: "/",
    icon: "material-symbols:book",
  },
] as const;

export const SidebarNavigation: React.FC = () => {
  const pathname = usePathname();

  return (
    <aside
      className={[
        "hidden",
        "md:fixed",
        "md:left-0",
        "md:top-0",
        "md:flex",
        "md:h-screen",
        "md:w-64",
        "md:flex-col",
        "md:border-r",
        "md:border-slate-200",
        "md:bg-white",
        "md:dark:border-slate-800",
        "md:dark:bg-slate-900",
      ].join(" ")}
    >
      <div className={["px-6", "py-6"].join(" ")}>
        <Link href="/">
          <Image
            src={holocronLogo}
            alt="HOLOCRON"
            className={["h-8", "w-auto", "dark:invert"].join(" ")}
          />
        </Link>
      </div>

      <nav className={["flex-1", "px-3"].join(" ")}>
        <ul className={["space-y-1"].join(" ")}>
          {navItems.map((item, index) => {
            const isActive = pathname === item.href;
            return (
              <Fragment key={item.href}>
                <li>
                  <Link
                    href={item.href}
                    aria-current={isActive ? "page" : undefined}
                    className={[
                      "flex",
                      "items-center",
                      "gap-4",
                      "rounded-lg",
                      "px-4",
                      "py-3",
                      "text-base",
                      "font-semibold",
                      "transition-colors",
                      isActive
                        ? "bg-sky-50 text-sky-700 dark:bg-sky-950 dark:text-sky-500"
                        : "text-slate-700 hover:bg-slate-50 dark:text-slate-300 dark:hover:bg-slate-800",
                    ].join(" ")}
                  >
                    <Icon icon={item.icon} className="size-6" />
                    <span>{item.label}</span>
                  </Link>
                </li>
                {/* ライブラリの後に登録ボタンを挿入 */}
                {index === 0 && (
                  <li key="add-book">
                    <Link
                      href="/books/new"
                      className={[
                        "flex",
                        "items-center",
                        "gap-4",
                        "rounded-lg",
                        "px-4",
                        "py-3",
                        "text-base",
                        "font-semibold",
                        "text-slate-700",
                        "transition-colors",
                        "hover:bg-slate-50",
                        "dark:text-slate-300",
                        "dark:hover:bg-slate-800",
                      ].join(" ")}
                    >
                      <Icon icon="material-symbols:add" className="size-6" />
                      <span>書籍を登録</span>
                    </Link>
                  </li>
                )}
              </Fragment>
            );
          })}
        </ul>
      </nav>
    </aside>
  );
};
