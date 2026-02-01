import { Icon } from "@iconify/react";
import type { FC } from "react";
import type { RegisteredBook } from "../page.view";

export const RegisteredBookCard: FC<{ book: RegisteredBook }> = ({ book }) => (
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
    <div
      className={[
        "flex",
        "items-center",
        "gap-2",
        "mb-3",
        "text-slate-900",
        "dark:text-slate-100",
      ].join(" ")}
    >
      <Icon icon="material-symbols:check-circle" className="size-5" />
      <span className="text-sm font-medium">登録しました</span>
    </div>

    <div className={["flex", "gap-4"].join(" ")}>
      <div
        className={[
          "h-28",
          "w-20",
          "shrink-0",
          "overflow-hidden",
          "rounded-lg",
          "bg-slate-100",
          "dark:bg-slate-800",
        ].join(" ")}
      >
        <img
          src={book.thumbnailUrl}
          alt={book.title}
          className={["h-full", "w-full", "object-cover"].join(" ")}
        />
      </div>
      <div
        className={["flex", "flex-1", "flex-col", "justify-between"].join(" ")}
      >
        <div>
          <h2
            className={[
              "font-bold",
              "text-slate-900",
              "line-clamp-2",
              "dark:text-slate-100",
            ].join(" ")}
          >
            {book.title}
          </h2>
          <p
            className={[
              "mt-1",
              "text-sm",
              "text-slate-500",
              "dark:text-slate-400",
            ].join(" ")}
          >
            {book.authors.join(", ")}
          </p>
        </div>
      </div>
    </div>
  </div>
);
