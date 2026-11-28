import type { APIRoute } from 'astro';
import guide from '../content/docs/index.md?raw';

export const GET: APIRoute = () => new Response(guide, {
  headers: { 'Content-Type': 'text/markdown; charset=utf-8' },
});
