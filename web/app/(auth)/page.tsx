"use client";

import { useSuspenseQuery } from "@tanstack/react-query";
import { Suspense, useDeferredValue, useState } from "react";
import { useAuth } from "../_components/AuthProvider";
import { fetchClient } from "../_lib/query";
import type { Book } from "./_models/book";
import { HomePage } from "./page.view";

function HomeContent({
  query,
  deferredQuery,
  onChangeQuery,
}: {
  query: string;
  deferredQuery: string;
  onChangeQuery: (query: string) => void;
}) {
  const authState = useAuth();

  if (authState.status !== "authenticated") {
    throw new Error("HomeContent requires authenticated user");
  }

  const { user } = authState;

  const { data } = useSuspenseQuery({
    queryKey: ["books", deferredQuery, user.uid],
    queryFn: async () => {
      const token = await user.getIdToken();
      const { data } = await fetchClient.GET("/books", {
        params: { query: deferredQuery ? { q: deferredQuery } : {} },
        headers: { Authorization: `Bearer ${token}` },
      });
      return data ?? { items: [] };
    },
  });

  const books = (data?.items ?? []) as Book[];

  return <HomePage books={books} query={query} onChangeQuery={onChangeQuery} />;
}

export default function Home() {
  const authState = useAuth();
  const [query, setQuery] = useState("");
  const deferredQuery = useDeferredValue(query);

  if (authState.status !== "authenticated") {
    return <HomePage books={[]} query="" onChangeQuery={() => {}} loading />;
  }

  return (
    <Suspense
      fallback={
        <HomePage books={[]} query={query} onChangeQuery={setQuery} loading />
      }
    >
      <HomeContent
        query={query}
        deferredQuery={deferredQuery}
        onChangeQuery={setQuery}
      />
    </Suspense>
  );
}
