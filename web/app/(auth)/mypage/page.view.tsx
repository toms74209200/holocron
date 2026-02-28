"use client";

import { Icon } from "@iconify/react";
import Link from "next/link";
import type { FC } from "react";
import { type Borrowing, isOverdue } from "../_models/borrowing";

interface MyPageProps {
  borrowings: Borrowing[];
  isLoading?: boolean;
}

export const MyPage: FC<MyPageProps> = ({ borrowings, isLoading }) => {
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
            "px-4",
            "py-4",
          ].join(" ")}
        >
          <h1
            className={[
              "text-xl",
              "font-bold",
              "text-slate-900",
              "dark:text-slate-100",
            ].join(" ")}
          >
            マイページ
          </h1>
        </div>
      </header>

      <main className={["mx-auto", "max-w-2xl", "p-4", "space-y-4"].join(" ")}>
        <section>
          <h2
            className={[
              "mb-3",
              "text-base",
              "font-semibold",
              "text-slate-700",
              "dark:text-slate-300",
            ].join(" ")}
          >
            借りている本
          </h2>

          {isLoading ? (
            <LoadingState />
          ) : borrowings.length === 0 ? (
            <EmptyState />
          ) : (
            <ul className={["space-y-3"].join(" ")}>
              {borrowings.map((borrowing) => (
                <li key={borrowing.id}>
                  <BorrowingCard borrowing={borrowing} />
                </li>
              ))}
            </ul>
          )}
        </section>
      </main>
    </div>
  );
};

const BorrowingCard: FC<{ borrowing: Borrowing }> = ({ borrowing }) => {
  const borrowedAt = new Date(borrowing.borrowedAt).toLocaleDateString("ja-JP");
  const dueDate = borrowing.dueDate
    ? new Date(borrowing.dueDate).toLocaleDateString("ja-JP")
    : null;
  const overdue = isOverdue(borrowing);

  return (
    <Link
      href={`/book?id=${borrowing.id}`}
      className={[
        "flex",
        "gap-4",
        "rounded-xl",
        "border",
        "border-slate-200",
        "bg-white",
        "p-4",
        "transition-colors",
        "hover:bg-slate-50",
        "dark:border-slate-800",
        "dark:bg-slate-900",
        "dark:hover:bg-slate-800",
      ].join(" ")}
    >
      <div
        className={[
          "h-24",
          "w-16",
          "shrink-0",
          "overflow-hidden",
          "rounded-lg",
          "bg-slate-100",
          "dark:bg-slate-800",
        ].join(" ")}
      >
        {borrowing.thumbnailUrl ? (
          // biome-ignore lint/performance/noImgElement: サムネイルは外部URL
          <img
            src={borrowing.thumbnailUrl}
            alt={borrowing.title}
            className={["h-full", "w-full", "object-cover"].join(" ")}
          />
        ) : (
          <div
            className={[
              "flex",
              "h-full",
              "w-full",
              "items-center",
              "justify-center",
            ].join(" ")}
          >
            <Icon
              icon="material-symbols:book"
              className={["size-8", "text-slate-400"].join(" ")}
            />
          </div>
        )}
      </div>

      <div
        className={["flex", "flex-1", "flex-col", "justify-between"].join(" ")}
      >
        <div>
          <p
            className={[
              "font-bold",
              "text-slate-900",
              "line-clamp-2",
              "dark:text-slate-100",
            ].join(" ")}
          >
            {borrowing.title}
          </p>
          <p
            className={[
              "mt-1",
              "text-sm",
              "text-slate-500",
              "dark:text-slate-400",
            ].join(" ")}
          >
            {borrowing.authors.join(", ")}
          </p>
        </div>

        <div className={["mt-2", "space-y-1"].join(" ")}>
          <p
            className={[
              "text-xs",
              "text-slate-500",
              "dark:text-slate-400",
            ].join(" ")}
          >
            貸出日: {borrowedAt}
          </p>
          {dueDate && (
            <p
              className={[
                "text-xs",
                "font-medium",
                overdue
                  ? "text-red-600 dark:text-red-400"
                  : "text-slate-500 dark:text-slate-400",
              ].join(" ")}
            >
              返却期限: {dueDate}
              {overdue && " (延滞中)"}
            </p>
          )}
        </div>
      </div>

      <div className={["flex", "items-center", "self-center"].join(" ")}>
        <Icon
          icon="material-symbols:chevron-right"
          className={["size-5", "text-slate-400"].join(" ")}
        />
      </div>
    </Link>
  );
};

const LoadingState: FC = () => (
  <div
    className={[
      "flex",
      "flex-col",
      "items-center",
      "justify-center",
      "rounded-xl",
      "border",
      "border-slate-200",
      "bg-white",
      "p-12",
      "dark:border-slate-800",
      "dark:bg-slate-900",
    ].join(" ")}
  >
    <Icon
      icon="material-symbols:progress-activity"
      className={["size-8", "animate-spin", "text-slate-400"].join(" ")}
    />
  </div>
);

const EmptyState: FC = () => (
  <div
    className={[
      "flex",
      "flex-col",
      "items-center",
      "justify-center",
      "rounded-xl",
      "border",
      "border-slate-200",
      "bg-white",
      "p-12",
      "text-center",
      "dark:border-slate-800",
      "dark:bg-slate-900",
    ].join(" ")}
  >
    <Icon
      icon="material-symbols:book-outline"
      className={[
        "mb-3",
        "size-12",
        "text-slate-300",
        "dark:text-slate-600",
      ].join(" ")}
    />
    <p
      className={["text-sm", "text-slate-500", "dark:text-slate-400"].join(" ")}
    >
      現在借りている本はありません
    </p>
  </div>
);
