import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import { unified } from '@astrojs/markdown-remark';
import { tableScroll } from '@cloudflare/nimbus-docs/markdown';
import { siteBase, siteOrigin } from './site-base.mjs';

// GitHub project sites live under /<repository>/; local previews default to /.
const base = siteBase();

export default defineConfig({
  site: siteOrigin(),
  base,
  vite: {
    server: {
      allowedHosts: ['.e2b.app'],
      proxy: { [`${base}reference`]: { target: 'http://127.0.0.1:4322', changeOrigin: true, ws: true } },
    },
    preview: { allowedHosts: ['.e2b.app'] },
  },
  markdown: {
    // Nimbus and Starlight are separate frameworks. Reuse Nimbus's accessible
    // table wrapper in Starlight's unified pipeline instead of running two
    // conflicting site integrations (and two search/sitemap builders).
    processor: unified({ rehypePlugins: [function nimbusAccessibleTables() {
      const wrap = tableScroll().element.visit;
      return (tree) => {
        const visit = (node) => {
          if (!node.children) return;
          node.children = node.children.map((child) => {
            if (child.type === 'element' && child.tagName === 'table') return wrap(child) || child;
            visit(child);
            return child;
          });
        };
        visit(tree);
      };
    }] }),
  },
  integrations: [starlight({
    title: 'FCC Docs',
    defaultLocale: 'root',
    locales: {
      root: { label: 'English', lang: 'en' },
      vi: { label: 'Tiếng Việt', lang: 'vi' },
    },
    social: [{ icon: 'github', label: 'GitHub', href: 'https://github.com/MoLik-Labs/FCC' }],
    customCss: ['./src/styles/custom.css'],
    expressiveCode: {
      plugins: [{
        name: 'Keyboard-accessible code blocks',
        hooks: { postprocessRenderedBlock({ renderData }) {
          const visit = (node) => {
            if (node.type === 'element' && node.tagName === 'pre') node.properties.tabIndex = 0;
            for (const child of node.children || []) visit(child);
          };
          visit(renderData.blockAst);
        } },
      }],
    },
    sidebar: [{
      label: 'Field guide', translations: { vi: 'Hướng dẫn' },
      items: [
        { label: 'Overview', translations: { vi: 'Tổng quan' }, slug: '' },
        ...[
          ['quickstart', 'Quick start', 'Cài đặt nhanh'],
          ['first-request', 'First request', 'Yêu cầu đầu tiên'],
          ['providers', 'Providers and models', 'Nhà cung cấp và mô hình'],
          ['clients', 'Clients', 'Công cụ kết nối'],
          ['api', 'API', 'API'],
          ['messaging', 'Messaging', 'Nhắn tin'],
          ['workloads', 'Workloads', 'Tải công việc'],
          ['security', 'Security', 'Bảo mật'],
          ['maintenance', 'Maintenance', 'Bảo trì'],
          ['troubleshooting', 'Troubleshooting', 'Khắc phục sự cố'],
          ['development', 'Development', 'Phát triển'],
          ['pages', 'Publishing', 'Xuất bản'],
          ['resources', 'Resources', 'Tài nguyên'],
        ].map(([id, label, vi]) => ({ label, translations: { vi }, link: `#${id}` })),
      ],
    }],
  })],
});
