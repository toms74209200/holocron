import { Temporal } from "temporal-polyfill";
import { expect, test } from "vitest";
import { type Borrowing, isOverdue } from "./borrowing";

test("when isOverdue with no dueDate then returns false", () => {
  const borrowing: Borrowing = {
    id: crypto.randomUUID(),
    title: crypto.randomUUID(),
    authors: [crypto.randomUUID()],
    borrowedAt: Temporal.Now.instant().toString(),
  };

  expect(isOverdue(borrowing)).toBe(false);
});

test("when isOverdue with future dueDate then returns false", () => {
  const borrowing: Borrowing = {
    id: crypto.randomUUID(),
    title: crypto.randomUUID(),
    authors: [crypto.randomUUID()],
    borrowedAt: Temporal.Now.instant().subtract({ hours: 48 }).toString(),
    dueDate: Temporal.Now.instant().add({ hours: 24 }).toString(),
  };

  expect(isOverdue(borrowing)).toBe(false);
});

test("when isOverdue with past dueDate then returns true", () => {
  const borrowing: Borrowing = {
    id: crypto.randomUUID(),
    title: crypto.randomUUID(),
    authors: [crypto.randomUUID()],
    borrowedAt: Temporal.Now.instant().subtract({ hours: 48 }).toString(),
    dueDate: Temporal.Now.instant().subtract({ hours: 24 }).toString(),
  };

  expect(isOverdue(borrowing)).toBe(true);
});
