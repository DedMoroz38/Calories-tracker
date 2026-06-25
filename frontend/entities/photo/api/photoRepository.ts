/**
 * entities/photo/api/photoRepository.ts
 *
 * Photo posts + public feed. Images live in S3; the API returns short-lived
 * presigned URLs, so treat `url` fields as ephemeral (don't cache long-term).
 */

import { http, uploadFile } from "@/shared/api/httpClient";

export interface Photo {
  id: number;
  url: string;
  created_at: string;
}

export interface FeedItem {
  id: number;
  url: string;
  user_id: number;
  author_name: string;
  author_avatar: string;
  created_at: string;
}

export interface FeedPage {
  items: FeedItem[];
  next_cursor: number;
}

/** Post a new photo to your profile. */
export async function uploadPhoto(file: File): Promise<Photo> {
  return uploadFile<Photo>("/photos", file);
}

/** Your own photos, newest first. */
export async function getMyPhotos(): Promise<Photo[]> {
  return http.get<Photo[]>("/photos/me");
}

/** Delete one of your photos. */
export async function deletePhoto(id: number): Promise<void> {
  await http.delete(`/photos/${id}`);
}

/** A page of the public feed (other users' photos). */
export async function getFeed(cursor?: number, limit = 12): Promise<FeedPage> {
  const params = new URLSearchParams({ limit: String(limit) });
  if (cursor) params.set("cursor", String(cursor));
  return http.get<FeedPage>(`/feed?${params.toString()}`);
}

/** Replace your avatar; returns the new presigned avatar URL. */
export async function uploadAvatar(file: File): Promise<{ avatar_url: string }> {
  return uploadFile<{ avatar_url: string }>("/profile/avatar", file);
}
