// Keep stable historic URLs without maintaining duplicate copies of the old pages.
import { copyFile, mkdir, readFile, rm, writeFile } from 'node:fs/promises';
import { resolve } from 'node:path';

await mkdir('public/assets', { recursive: true });
for (const name of ['style.css', 'site.js', 'molik.css', 'molik.html']) {
  await copyFile(name, `public/${name}`);
}
for (const [source, target, lang] of [['index.html', 'legacy.html', 'en'], ['vi.html', 'vi.html', 'vi']]) {
  let html = await readFile(source, 'utf8');
  const notice = lang === 'en'
    ? 'This is the previous guide. <a href="index.html">Open the current Starlight guide</a>. Publishing instructions on this archived page are outdated; see the current guide.'
    : 'Đây là hướng dẫn trước đây. <a href="vi/">Mở hướng dẫn Starlight hiện tại</a>. Hướng dẫn xuất bản trong trang cũ không còn cập nhật; xem hướng dẫn mới.';
  html = html.replace('<body>', `<body data-pagefind-ignore><div style="padding:1rem;background:#dcefe3;color:#123f29">${notice}</div>`);
  await writeFile(`public/${target}`, html);
}
await copyFile(resolve('../../ohmymolik/brand/molik-mark.svg'), 'public/assets/molik-mark.svg');

// Astro caches rendered Markdown in node_modules; invalidate it so changes to
// rendering and accessibility plugins are reflected in every build.
await rm('node_modules/.astro/data-store.json', { force: true });
