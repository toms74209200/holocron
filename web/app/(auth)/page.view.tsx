"use client";

import { Icon } from "@iconify/react";
import Link from "next/link";
import { BookInfoCard } from "./_components/BookInfoCard";
import type { Book } from "./_models/book";

interface HomePageProps {
  books: Book[];
  query: string;
  onChangeQuery: (query: string) => void;
  loading?: boolean;
}

export const HomePage: React.FC<HomePageProps> = ({
  books,
  query,
  onChangeQuery,
  loading,
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
            "justify-between",
            "px-4",
            "py-4",
          ].join(" ")}
        >
          <h1
            className={[
              "text-xl",
              "font-bold",
              "tracking-widest",
              "text-slate-900",
              "dark:text-slate-100",
            ].join(" ")}
          >
            HOLOCRON
          </h1>
          <Link
            href="/books/new"
            className={[
              "flex",
              "items-center",
              "gap-1",
              "rounded-lg",
              "bg-sky-700",
              "px-4",
              "py-2",
              "text-sm",
              "font-semibold",
              "text-white",
              "transition-colors",
              "hover:bg-sky-800",
              "dark:bg-sky-600",
              "dark:hover:bg-sky-700",
            ].join(" ")}
          >
            <Icon icon="material-symbols:add" className="h-4 w-4" />
            登録
          </Link>
        </div>
      </header>

      <main className={["mx-auto", "max-w-4xl", "p-4"].join(" ")}>
        <div className={["mb-6"].join(" ")}>
          <div className={["relative"].join(" ")}>
            <label
              htmlFor="search"
              className={[
                "absolute",
                "left-4",
                "top-1/2",
                "-translate-y-1/2",
                "cursor-text",
              ].join(" ")}
            >
              <Icon
                icon="material-symbols:search"
                className={["h-5", "w-5", "text-slate-400"].join(" ")}
              />
            </label>
            <input
              id="search"
              type="text"
              placeholder="タイトル・著者で検索..."
              value={query}
              onChange={(e) => onChangeQuery(e.target.value)}
              className={[
                "w-full",
                "rounded-lg",
                "border",
                "border-slate-200",
                "bg-white",
                "py-3",
                "pl-12",
                "pr-4",
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
        </div>

        {loading ? (
          <div
            className={["flex", "items-center", "justify-center", "py-16"].join(
              " ",
            )}
          >
            <p
              className={[
                "text-sm",
                "text-slate-500",
                "dark:text-slate-400",
              ].join(" ")}
            >
              読み込み中...
            </p>
          </div>
        ) : books.length === 0 ? (
          <div
            className={[
              "flex",
              "flex-col",
              "items-center",
              "justify-center",
              "py-16",
              "text-center",
            ].join(" ")}
          >
            <p className={["text-slate-500", "dark:text-slate-400"].join(" ")}>
              書籍がありません
            </p>
            <p
              className={[
                "mt-1",
                "text-sm",
                "text-slate-400",
                "dark:text-slate-500",
              ].join(" ")}
            >
              右上の「+ 登録」から書籍を追加してください
            </p>
          </div>
        ) : (
          <ul className={["space-y-3"].join(" ")}>
            {books.map((book) => (
              <li key={book.id}>
                <Link
                  href={`/books/${book.id}`}
                  className={[
                    "block",
                    "rounded-xl",
                    "border",
                    "border-slate-200",
                    "bg-white",
                    "p-4",
                    "transition-colors",
                    "hover:border-slate-300",
                    "dark:border-slate-800",
                    "dark:bg-slate-900",
                    "dark:hover:border-slate-700",
                  ].join(" ")}
                >
                  <BookInfoCard
                    title={book.title}
                    authors={book.authors}
                    status={book.status}
                    borrower={book.borrower}
                    thumbnailUrl={book.thumbnailUrl}
                  />
                </Link>
              </li>
            ))}
          </ul>
        )}
      </main>
    </div>
  );
};
