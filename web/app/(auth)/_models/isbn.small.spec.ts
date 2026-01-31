import { expect, test } from "vitest";
import { parseIsbn } from "./isbn";

test.each([
  {
    name: "when parseIsbn with ISBN-13 hyphenated then returns normalized ISBN-13",
    input: "978-0-306-40615-7",
    expected: "9780306406157",
  },
  {
    name: "when parseIsbn with ISBN-10 X check digit then returns normalized ISBN-10",
    input: "0-8044-2957-X",
    expected: "080442957X",
  },
  {
    name: "when parseIsbn with invalid checksum then returns null",
    input: "9780306406158",
    expected: null,
  },
  {
    name: "when parseIsbn with non numeric characters then returns null",
    input: "ABC-DEF",
    expected: null,
  },
])("$name", ({ input, expected }) => {
  const result = parseIsbn(input);
  expect(result).toBe(expected);
});
