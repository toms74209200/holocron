"use client";

import { Icon } from "@iconify/react";
import Link from "next/link";
import type { FC } from "react";

interface NewBookPageProps {
  code: string;
  error?: string | null;
  onChangeCode: (code: string) => void;
  onSubmit: (e: { preventDefault: () => void }) => void;
}

export const NewBookPage: FC<NewBookPageProps> = ({
  code,
  error,
  onChangeCode,
  onSubmit,
}) => {
  return (
    <div
      className={["min-h-screen", "bg-slate-50", "dark:bg-slate-950"].join(" ")}
    >
      <header
        className={[
          "border-b",
          "border-slate-200",
          "bg-white",
          "dark:border-slate-800",
          "dark:bg-slate-900",
        ].join(" ")}
      >
        <div
          className={[
            "mx-auto",
            "flex",
            "max-w-4xl",
            "items-center",
            "gap-4",
            "px-4",
            "py-4",
          ].join(" ")}
        >
          <Link
            href="/"
            className={[
              "flex",
              "items-center",
              "gap-1",
              "text-slate-600",
              "transition-colors",
              "hover:text-slate-900",
              "dark:text-slate-400",
              "dark:hover:text-slate-100",
            ].join(" ")}
          >
            <Icon
              icon="material-symbols:arrow-back"
              className={["h-4", "w-4"].join(" ")}
            />
            戻る
          </Link>
          <h1
            className={[
              "text-xl",
              "font-bold",
              "text-slate-900",
              "dark:text-slate-100",
            ].join(" ")}
          >
            書籍登録
          </h1>
        </div>
      </header>

      <main className={["mx-auto", "max-w-md", "p-4"].join(" ")}>
        <div
          className={[
            "rounded-xl",
            "border",
            "border-slate-200",
            "bg-white",
            "p-6",
            "dark:border-slate-800",
            "dark:bg-slate-900",
          ].join(" ")}
        >
          <form onSubmit={onSubmit} className={["space-y-4"].join(" ")}>
            <div>
              <label
                htmlFor="code"
                className={[
                  "mb-2",
                  "block",
                  "text-sm",
                  "font-medium",
                  "text-slate-700",
                  "dark:text-slate-300",
                ].join(" ")}
              >
                ISBN
              </label>
              <input
                id="code"
                type="text"
                value={code}
                onChange={(e) => onChangeCode(e.target.value)}
                placeholder="978-4-87311-752-2"
                className={[
                  "w-full",
                  "rounded-lg",
                  "border",
                  "border-slate-200",
                  "bg-white",
                  "px-4",
                  "py-3",
                  "text-slate-900",
                  "placeholder:text-slate-400",
                  "focus:border-sky-600",
                  "focus:outline-none",
                  "focus:ring-1",
                  "focus:ring-sky-600",
                  "dark:border-slate-700",
                  "dark:bg-slate-800",
                  "dark:text-slate-100",
                  "dark:placeholder:text-slate-500",
                  "dark:focus:border-sky-500",
                  "dark:focus:ring-sky-500",
                ].join(" ")}
              />
            </div>

            {error && (
              <div
                className={[
                  "rounded-lg",
                  "bg-red-50",
                  "px-4",
                  "py-3",
                  "text-sm",
                  "text-red-600",
                  "dark:bg-red-500/10",
                  "dark:text-red-400",
                ].join(" ")}
              >
                {error}
              </div>
            )}

            <button
              type="submit"
              disabled={!code.trim()}
              className={[
                "flex",
                "w-full",
                "items-center",
                "justify-center",
                "gap-1",
                "rounded-lg",
                "bg-sky-700",
                "px-4",
                "py-3",
                "font-semibold",
                "text-white",
                "transition-colors",
                "hover:bg-sky-800",
                "disabled:cursor-not-allowed",
                "disabled:bg-slate-300",
                "disabled:text-slate-500",
                "dark:bg-sky-600",
                "dark:hover:bg-sky-700",
                "dark:disabled:bg-slate-700",
                "dark:disabled:text-slate-400",
              ].join(" ")}
            >
              <Icon
                icon="material-symbols:add"
                className={["h-5", "w-5"].join(" ")}
              />
              登録する
            </button>
          </form>
        </div>
      </main>
    </div>
  );
};
