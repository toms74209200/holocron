"use client";

import { Icon } from "@iconify/react";
import Link from "next/link";
import type { FC } from "react";
import { Temporal } from "temporal-polyfill";
import { BookInfoCard } from "../../_components/BookInfoCard";
import type { Book } from "../../_models/book";

type ActionState =
  | { status: "idle" }
  | { status: "selecting-date" }
  | { status: "loading"; action: "borrow" | "return" }
  | { status: "error"; message: string };

interface BookDetailPageProps {
  book: Book;
  actionState: ActionState;
  currentUserId: string;
  dueDate: string;
  onBorrowClick: () => void;
  onBorrowConfirm: () => void;
  onBorrowCancel: () => void;
  onDueDateChange: (date: string) => void;
  onReturn: () => void;
}

export const BookDetailPage: FC<BookDetailPageProps> = ({
  book,
  actionState,
  currentUserId,
  dueDate,
  onBorrowClick,
  onBorrowConfirm,
  onBorrowCancel,
  onDueDateChange,
  onReturn,
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
            書籍詳細
          </h1>
        </div>
      </header>

      <main className={["mx-auto", "max-w-md", "p-4", "space-y-4"].join(" ")}>
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
          <BookInfoCard
            title={book.title}
            authors={book.authors}
            status={book.status}
            borrower={book.borrower}
            thumbnailUrl={book.thumbnailUrl}
          />

          {(book.publisher || book.publishedDate || book.code) && (
            <div
              className={[
                "mt-4",
                "pt-4",
                "border-t",
                "border-slate-200",
                "dark:border-slate-800",
                "space-y-2",
              ].join(" ")}
            >
              {book.publisher && (
                <div>
                  <p
                    className={[
                      "text-xs",
                      "text-slate-500",
                      "dark:text-slate-400",
                    ].join(" ")}
                  >
                    出版社
                  </p>
                  <p
                    className={[
                      "text-sm",
                      "text-slate-900",
                      "dark:text-slate-100",
                    ].join(" ")}
                  >
                    {book.publisher}
                  </p>
                </div>
              )}
              {book.publishedDate && (
                <div>
                  <p
                    className={[
                      "text-xs",
                      "text-slate-500",
                      "dark:text-slate-400",
                    ].join(" ")}
                  >
                    出版日
                  </p>
                  <p
                    className={[
                      "text-sm",
                      "text-slate-900",
                      "dark:text-slate-100",
                    ].join(" ")}
                  >
                    {book.publishedDate}
                  </p>
                </div>
              )}
              {book.code && (
                <div>
                  <p
                    className={[
                      "text-xs",
                      "text-slate-500",
                      "dark:text-slate-400",
                    ].join(" ")}
                  >
                    コード
                  </p>
                  <p
                    className={[
                      "text-sm",
                      "text-slate-900",
                      "dark:text-slate-100",
                    ].join(" ")}
                  >
                    {book.code}
                  </p>
                </div>
              )}
            </div>
          )}
        </div>

        {actionState.status === "error" && (
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
            {actionState.message}
          </div>
        )}

        {book.status === "available" &&
          (actionState.status === "selecting-date" ? (
            <div
              className={[
                "rounded-xl",
                "border",
                "border-slate-200",
                "bg-white",
                "p-6",
                "dark:border-slate-800",
                "dark:bg-slate-900",
                "space-y-4",
              ].join(" ")}
            >
              <div>
                <label
                  htmlFor="dueDate"
                  className={[
                    "mb-2",
                    "block",
                    "text-sm",
                    "font-medium",
                    "text-slate-700",
                    "dark:text-slate-300",
                  ].join(" ")}
                >
                  返却期限
                </label>
                <input
                  id="dueDate"
                  type="date"
                  value={dueDate}
                  onChange={(e) => onDueDateChange(e.target.value)}
                  min={Temporal.Now.plainDateISO().toString()}
                  className={[
                    "w-full",
                    "rounded-lg",
                    "border",
                    "border-slate-200",
                    "bg-white",
                    "px-4",
                    "py-3",
                    "text-slate-900",
                    "focus:border-sky-600",
                    "focus:outline-none",
                    "focus:ring-1",
                    "focus:ring-sky-600",
                    "dark:border-slate-700",
                    "dark:bg-slate-800",
                    "dark:text-slate-100",
                    "dark:focus:border-sky-500",
                    "dark:focus:ring-sky-500",
                  ].join(" ")}
                />
              </div>

              <div className={["flex", "gap-2"].join(" ")}>
                <button
                  type="button"
                  onClick={onBorrowCancel}
                  className={[
                    "flex",
                    "flex-1",
                    "items-center",
                    "justify-center",
                    "gap-1",
                    "rounded-lg",
                    "border",
                    "border-slate-200",
                    "bg-white",
                    "px-4",
                    "py-3",
                    "font-semibold",
                    "text-slate-900",
                    "transition-colors",
                    "hover:bg-slate-50",
                    "dark:border-slate-700",
                    "dark:bg-slate-800",
                    "dark:text-slate-100",
                    "dark:hover:bg-slate-700",
                  ].join(" ")}
                >
                  キャンセル
                </button>
                <button
                  type="button"
                  onClick={onBorrowConfirm}
                  disabled={!dueDate}
                  className={[
                    "flex",
                    "flex-1",
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
                  確定
                </button>
              </div>
            </div>
          ) : (
            <button
              type="button"
              onClick={onBorrowClick}
              disabled={actionState.status === "loading"}
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
              {actionState.status === "loading" &&
              actionState.action === "borrow" ? (
                <>
                  <Icon
                    icon="material-symbols:progress-activity"
                    className="h-5 w-5 animate-spin"
                  />
                  借りています...
                </>
              ) : (
                <>
                  <Icon icon="material-symbols:book" className="h-5 w-5" />
                  借りる
                </>
              )}
            </button>
          ))}

        {book.status === "borrowed" && book.borrower?.id === currentUserId && (
          <button
            type="button"
            onClick={onReturn}
            disabled={actionState.status === "loading"}
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
            {actionState.status === "loading" &&
            actionState.action === "return" ? (
              <>
                <Icon
                  icon="material-symbols:progress-activity"
                  className="h-5 w-5 animate-spin"
                />
                返却中...
              </>
            ) : (
              <>
                <Icon
                  icon="material-symbols:keyboard-return"
                  className="h-5 w-5"
                />
                返却する
              </>
            )}
          </button>
        )}
      </main>
    </div>
  );
};
