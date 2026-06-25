"use client";

import { useState, useEffect, useCallback, useRef } from "react";
import { TabBar } from "@/shared/ui";
import { RouteGuard } from "@/features/auth/ui/RouteGuard";
import { getFeed, uploadPhoto, type FeedItem } from "@/entities/photo/api/photoRepository";
import styles from "./feed.module.css";

function timeAgo(iso: string): string {
  const diff = Date.now() - new Date(iso).getTime();
  const m = Math.floor(diff / 60000);
  if (m < 1) return "just now";
  if (m < 60) return `${m}m ago`;
  const h = Math.floor(m / 60);
  if (h < 24) return `${h}h ago`;
  const d = Math.floor(h / 24);
  return `${d}d ago`;
}

function FeedScreen() {
  const [items, setItems]     = useState<FeedItem[]>([]);
  const [cursor, setCursor]   = useState<number | undefined>(undefined);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);
  const [err, setErr]         = useState<string | null>(null);
  const [toast, setToast]     = useState<string | null>(null);

  const photoInput = useRef<HTMLInputElement>(null);
  const sentinel   = useRef<HTMLDivElement>(null);
  // Guards against the observer firing again mid-request.
  const loadingRef = useRef(false);

  const loadMore = useCallback(async () => {
    if (loadingRef.current || !hasMore) return;
    loadingRef.current = true;
    setLoading(true);
    setErr(null);
    try {
      const page = await getFeed(cursor);
      setItems((prev) => [...prev, ...page.items]);
      setHasMore(page.next_cursor !== 0);
      setCursor(page.next_cursor || undefined);
    } catch (e) {
      setErr((e as Error).message);
      setHasMore(false);
    } finally {
      loadingRef.current = false;
      setLoading(false);
    }
  }, [cursor, hasMore]);

  // Initial load.
  useEffect(() => { loadMore(); /* eslint-disable-next-line react-hooks/exhaustive-deps */ }, []);

  // Infinite scroll: load the next page when the sentinel scrolls into view.
  useEffect(() => {
    const el = sentinel.current;
    if (!el) return;
    const obs = new IntersectionObserver(
      (entries) => { if (entries[0].isIntersecting) loadMore(); },
      { rootMargin: "200px" },
    );
    obs.observe(el);
    return () => obs.disconnect();
  }, [loadMore]);

  const handlePost = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    e.target.value = "";
    if (!file) return;
    try {
      await uploadPhoto(file);
      setToast("Posted to your profile ✓");
      setTimeout(() => setToast(null), 2200);
    } catch (e) {
      setErr((e as Error).message);
    }
  };

  return (
    <>
      <main className={styles.root}>
        <h1 className={styles.title}>Feed</h1>

        {err && <p className={styles.errMsg}>{err}</p>}

        {items.map((it) => (
          <article key={it.id} className={styles.post}>
            <header className={styles.postHeader}>
              {it.author_avatar ? (
                // eslint-disable-next-line @next/next/no-img-element
                <img className={styles.authorAvatar} src={it.author_avatar} alt="" />
              ) : (
                <div className={styles.authorFallback}>{it.author_name.charAt(0).toUpperCase()}</div>
              )}
              <span className={styles.authorName}>{it.author_name}</span>
              <span className={styles.postTime}>{timeAgo(it.created_at)}</span>
            </header>
            {/* eslint-disable-next-line @next/next/no-img-element */}
            <img className={styles.postImg} src={it.url} alt="" loading="lazy" />
          </article>
        ))}

        {!loading && items.length === 0 && !err && (
          <p className={styles.empty}>No photos from other users yet.</p>
        )}

        {loading && <p className={styles.loadingMsg}>Loading…</p>}

        <div ref={sentinel} className={styles.sentinel} />
        <div style={{ height: 80 }} />
      </main>

      {toast && <div className={styles.toast}>{toast}</div>}

      <input ref={photoInput} type="file" accept="image/*" hidden onChange={handlePost} />
      <TabBar onFabClick={() => photoInput.current?.click()} fabLabel="Post a photo" />
    </>
  );
}

export function FeedPage() {
  return (
    <RouteGuard>
      <FeedScreen />
    </RouteGuard>
  );
}
