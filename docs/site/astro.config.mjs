import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import { unified } from '@astrojs/markdown-remark';

// GitHub project sites live under /<repository>/; local previews default to /.
const repository = process.env.GITHUB_REPOSITORY?.split('/')[1];
const base = process.env.DOCS_BASE || (process.env.GITHUB_ACTIONS && repository ? `/${repository}/` : '/');

export default defineConfig({
  site: process.env.DOCS_SITE || (process.env.GITHUB_REPOSITORY_OWNER
    ? `https://${process.env.GITHUB_REPOSITORY_OWNER.toLowerCase()}.github.io`
    : 'https://molik-labs.github.io'),
  base,
  vite: { server: { allowedHosts: ['.e2b.app'] }, preview: { allowedHosts: ['.e2b.app'] } },
  markdown: {
    processor: unified({ rehypePlugins: [function keyboardAccessibleTables() {
      return (tree) => {
        const visit = (node) => {
          if (node.type === 'element' && node.tagName === 'table') node.properties.tabIndex = 0;
          for (const child of node.children || []) visit(child);
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
