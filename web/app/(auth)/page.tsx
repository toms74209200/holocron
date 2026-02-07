"use client";

import { useInfiniteQuery } from "@tanstack/react-query";
import { useDeferredValue, useEffect, useRef, useState } from "react";
import { useAuth } from "../_components/AuthProvider";
import { fetchClient } from "../_lib/query";
import type { Book } from "./_models/book";
import { HomePage } from "./page.view";

const PAGE_SIZE = 20;

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
  const observerTarget = useRef<HTMLDivElement>(null);

  if (authState.status !== "authenticated") {
    throw new Error("HomeContent requires authenticated user");
  }

  const { user } = authState;

  const { data, fetchNextPage, hasNextPage, isFetchingNextPage, isLoading } =
    useInfiniteQuery({
      queryKey: ["books", deferredQuery, user.uid],
      queryFn: async ({ pageParam = 0 }) => {
        const token = await user.getIdToken();
        const { data } = await fetchClient.GET("/books", {
          params: {
            query: {
              ...(deferredQuery ? { q: deferredQuery } : {}),
              limit: PAGE_SIZE,
              offset: pageParam,
            },
          },
          headers: { Authorization: `Bearer ${token}` },
        });
        return data ?? { items: [], total: 0 };
      },
      getNextPageParam: (lastPage, allPages) => {
        const loadedCount = allPages.reduce(
          (acc, page) => acc + (page?.items?.length ?? 0),
          0,
        );
        const total = lastPage?.total ?? 0;
        return loadedCount < total ? loadedCount : undefined;
      },
      initialPageParam: 0,
    });

  const books =
    data?.pages.flatMap((page) => (page?.items ?? []) as Book[]) ?? [];

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0]?.isIntersecting && hasNextPage && !isFetchingNextPage) {
          fetchNextPage();
        }
      },
      { threshold: 0.1 },
    );

    const currentTarget = observerTarget.current;
    if (currentTarget) {
      observer.observe(currentTarget);
    }

    return () => {
      if (currentTarget) {
        observer.unobserve(currentTarget);
      }
      observer.disconnect();
    };
  }, [hasNextPage, isFetchingNextPage, fetchNextPage]);

  return (
    <HomePage
      books={books}
      query={query}
      onChangeQuery={onChangeQuery}
      loading={isLoading}
      observerTarget={observerTarget}
      isFetchingNextPage={isFetchingNextPage}
    />
  );
}

export default function Home() {
  const authState = useAuth();
  const [query, setQuery] = useState("");
  const deferredQuery = useDeferredValue(query);
  const observerTarget = useRef<HTMLDivElement>(null);

  if (authState.status !== "authenticated") {
    return (
      <HomePage
        books={[]}
        query=""
        onChangeQuery={() => {}}
        loading
        observerTarget={observerTarget}
        isFetchingNextPage={false}
      />
    );
  }

  return (
    <HomeContent
      query={query}
      deferredQuery={deferredQuery}
      onChangeQuery={setQuery}
    />
  );
}
