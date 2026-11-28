import { defineConfig } from 'vitepress';
import { siteBase, siteOrigin } from '../../site-base.mjs';

// The reference is a VitePress subsite of the Starlight GitHub Pages build.
const parentBase = siteBase();

export default defineConfig({
  title: 'FCC API',
  description: 'FCC gateway API reference',
  base: `${parentBase}reference/`,
  sitemap: { hostname: new URL(`${parentBase}reference/`, siteOrigin()).href },
  transformHead({ pageData }) {
    const page = pageData.relativePath.replace(/(^|\/)index\.md$/, '$1').replace(/\.md$/, '.html');
    return [['link', { rel: 'canonical', href: new URL(`${parentBase}reference/${page}`, siteOrigin()).href }]];
  },
  outDir: '../dist/reference',
  cleanUrls: false,
  vite: { server: { allowedHosts: ['.e2b.app'] } },
  // Starlight's Pagefind pass runs after both builds; mark the VitePress
  // content as indexable so the main site's search also finds API pages.
  transformHtml(html, id) {
    if (id.endsWith('404.html')) return html;
    if (!html.includes('<main class="main"')) throw new Error('VitePress main element changed: ' + id);
    return html.replace('<main class="main"', '<main data-pagefind-body class="main"');
  },
  locales: {
    root: {
      label: 'English', lang: 'en',
      themeConfig: { nav: [{ text: 'API reference', link: '/' }], sidebar: [{ text: 'Reference', items: [{ text: 'API overview', link: '/' }] }] },
    },
    vi: {
      label: 'Tiếng Việt', lang: 'vi',
      themeConfig: { nav: [{ text: 'Tham khảo API', link: '/vi/' }], sidebar: [{ text: 'Tham khảo', items: [{ text: 'Tổng quan API', link: '/vi/' }] }], outline: { label: 'Trong trang này' }, returnToTopLabel: 'Về đầu trang', sidebarMenuLabel: 'Danh mục' },
    },
  },
  themeConfig: {
    search: {
      provider: 'local',
      options: { locales: { vi: { translations: {
        button: { buttonText: 'Tìm kiếm', buttonAriaLabel: 'Tìm kiếm' },
        modal: {
          displayDetails: 'Hiện chi tiết', resetButtonTitle: 'Xóa tìm kiếm',
          backButtonTitle: 'Đóng tìm kiếm', noResultsText: 'Không tìm thấy kết quả',
          footer: { selectText: 'chọn', selectKeyAriaLabel: 'Enter', navigateText: 'điều hướng', navigateUpKeyAriaLabel: 'Mũi tên lên', navigateDownKeyAriaLabel: 'Mũi tên xuống', closeText: 'đóng', closeKeyAriaLabel: 'Escape' },
        },
      } } } },
    },
    socialLinks: [{ icon: 'github', link: 'https://github.com/MoLik-Labs/FCC' }],
    footer: { message: 'FCC gateway · Static API reference' },
  },
});
