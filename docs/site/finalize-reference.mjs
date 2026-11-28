// Publish agent-readable reference pages and include the VitePress sitemap.
import { copyFile, readFile, writeFile } from 'node:fs/promises';
import { siteBase, siteOrigin } from './site-base.mjs';

const indexPath = 'dist/sitemap-index.xml';
const index = await readFile(indexPath, 'utf8');
await readFile('dist/reference/sitemap.xml', 'utf8');
const url = new URL(`${siteBase()}reference/sitemap.xml`, siteOrigin()).href;
const safeUrl = url.replaceAll('&', '&amp;').replaceAll('<', '&lt;').replaceAll('>', '&gt;');
if (!index.includes('</sitemapindex>')) throw new Error('Starlight sitemap index is missing or changed.');
await writeFile(indexPath, index.replace('</sitemapindex>', `<sitemap><loc>${safeUrl}</loc></sitemap></sitemapindex>`));

for (const page of ['index.md', 'vi/index.md']) {
  await copyFile(`reference/${page}`, `dist/reference/${page}`);
}
