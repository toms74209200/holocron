"use client";

import { Icon } from "@iconify/react";
import type { FC } from "react";
import { Temporal } from "temporal-polyfill";

interface BorrowFormProps {
  isSelectingDate: boolean;
  isLoading: boolean;
  dueDate: string;
  onBorrowClick: () => void;
  onBorrowConfirm: () => void;
  onBorrowCancel: () => void;
  onDueDateChange: (date: string) => void;
}

export const BorrowForm: FC<BorrowFormProps> = ({
  isSelectingDate,
  isLoading,
  dueDate,
  onBorrowClick,
  onBorrowConfirm,
  onBorrowCancel,
  onDueDateChange,
}) => {
  if (isSelectingDate) {
    return (
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
    );
  }

  return (
    <button
      type="button"
      onClick={onBorrowClick}
      disabled={isLoading}
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
      {isLoading ? (
        <>
          <Icon
            icon="material-symbols:progress-activity"
            className="size-5 animate-spin"
          />
          借りています...
        </>
      ) : (
        <>
          <Icon icon="material-symbols:book" className="size-5" />
          借りる
        </>
      )}
    </button>
  );
};
