import type { APIRoute } from 'astro';

// Static, Nimbus-inspired discovery for agent readers. The base is essential
// for GitHub Pages project sites; no server or external service is required.
export const GET: APIRoute = ({ site }) => {
  if (!site) throw new Error('Set Astro site to publish the agent index.');
  const root = new URL(import.meta.env.BASE_URL, site);
  return new Response(`# FCC documentation\n\n> FCC is a gateway for coding tools and model providers. Python is the supported runtime; Go and Oh My Molik are experimental.\n\n## Guides\n\n- [English guide](${new URL('guide.md', root)}): Installation, configuration, security, API and operations.\n- [Hướng dẫn tiếng Việt](${new URL('vi/guide.md', root)}): Cài đặt, cấu hình, bảo mật, API và vận hành.\n\n## Other resources\n\n- [API reference in English](${new URL('reference/index.md', root)}): Static API overview.\n- [Tham khảo API tiếng Việt](${new URL('reference/vi/index.md', root)}): Tổng quan API.\n- [Oh My Molik](${new URL('molik.html', root)}): Experimental themes and read-only MCP setup.\n- [Repository](https://github.com/MoLik-Labs/FCC): Project source and status.\n`, {
    headers: { 'Content-Type': 'text/plain; charset=utf-8' },
  });
};
