import { Temporal } from "temporal-polyfill";

export type Borrowing = {
  id: string;
  code?: string;
  title: string;
  authors: string[];
  publisher?: string;
  publishedDate?: string;
  thumbnailUrl?: string;
  borrowedAt: string;
  dueDate?: string;
};

export const isOverdue = (borrowing: Borrowing): boolean => {
  if (!borrowing.dueDate) {
    return false;
  }
  const due = Temporal.Instant.from(borrowing.dueDate).toZonedDateTimeISO(
    Temporal.Now.timeZoneId(),
  );
  const now = Temporal.Now.zonedDateTimeISO();
  return Temporal.ZonedDateTime.compare(due, now) < 0;
};
