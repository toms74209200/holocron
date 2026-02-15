"use client";

import Image from "next/image";
import Link from "next/link";
import holocronLogo from "./holocron-logo.svg";

export const MobileHeader: React.FC = () => {
  return (
    <header
      className={[
        "sticky",
        "top-0",
        "z-40",
        "border-b",
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
          "px-4",
          "py-3",
        ].join(" ")}
      >
        <Link href="/">
          <Image
            src={holocronLogo}
            alt="HOLOCRON"
            className={["h-7", "w-auto", "dark:invert"].join(" ")}
          />
        </Link>
      </div>
    </header>
  );
};
