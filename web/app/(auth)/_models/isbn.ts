export type Isbn = string & { readonly __brand: unique symbol };

const checksumIsbn10 = (s: string) => {
  if (s.length !== 10) {
    return null;
  }
  return s.split("").reduce((acc, char, index) => {
    const digit = char === "X" ? 10 : Number(char);
    return acc + digit * (10 - index);
  }, 0) %
    11 ===
    0
    ? s
    : null;
};

const checksumIsbn13 = (s: string) => {
  if (s.length !== 13) {
    return null;
  }
  return s
    .split("")
    .map(Number)
    .reduce((acc, digit, index) => {
      const weight = index % 2 === 0 ? 1 : 3;
      return acc + digit * weight;
    }, 0) %
    10 ===
    0
    ? s
    : null;
};

export const parseIsbn = (value: string): Isbn | null => {
  const normalized = value.trim().replace(/[\s-]/g, "").toUpperCase();
  const result = [
    (input: string) => (input ? input : null),
    (input: string) => (/^[0-9X]+$/.test(input) ? input : null),
    (input: string) =>
      input.length === 10 || input.length === 13 ? input : null,
    (input: string) =>
      input.length === 10
        ? checksumIsbn10(input)
        : input.length === 13
          ? checksumIsbn13(input)
          : null,
  ].reduce<string | null>(
    (acc, guard) => (acc === null ? null : guard(acc)),
    normalized,
  );

  return result ? (result as Isbn) : null;
};
