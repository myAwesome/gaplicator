import type { Post, CreatePostInput } from '../types/post';

const BASE = '/api/posts';

export interface PaginatedPosts { data: Post[]; total: number; page: number; limit: number; }

export async function listPosts(page = 1, limit = 20, sortBy = 'id', sortDir: 'asc' | 'desc' = 'desc', filters: Record<string, string> = {}): Promise<PaginatedPosts> {
  const params = new URLSearchParams({ page: String(page), limit: String(limit), sort_by: sortBy, sort_dir: sortDir, ...filters });
  const res = await fetch(`${BASE}?${params}`, { headers: {  } });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function getPost(id: number): Promise<Post> {
  const res = await fetch(`${BASE}/${id}`, { headers: {  } });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function createPost(data: CreatePostInput): Promise<Post> {
  const res = await fetch(BASE, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function updatePost(id: number, data: Partial<CreatePostInput>): Promise<Post> {
  const res = await fetch(`${BASE}/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

export async function deletePost(id: number): Promise<void> {
  const res = await fetch(`${BASE}/${id}`, { method: 'DELETE' });
  if (!res.ok) throw new Error(await res.text());
}

export async function batchDeletePost(ids: number[]): Promise<void> {
  const res = await fetch(`${BASE}/batch`, {
    method: 'DELETE',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ ids }),
  });
  if (!res.ok) throw new Error(await res.text());
}
