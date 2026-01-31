export type Book = {
  id: string;
  title: string;
  authors: string[];
  status: "available" | "borrowed";
  borrower?: {
    id: string;
    name: string;
    borrowedAt: string;
  };
  thumbnailUrl: string;
};
