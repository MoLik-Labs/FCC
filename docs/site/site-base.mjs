// Shared by both generators so every link works under the same Pages prefix.
export function siteBase(env = process.env) {
  const repository = env.GITHUB_REPOSITORY?.split('/')[1];
  const configured = env.DOCS_BASE || (env.GITHUB_ACTIONS && repository ? `/${repository}/` : '/');
  const segments = configured.split('/').filter(Boolean);
  return segments.length ? `/${segments.join('/')}/` : '/';
}

export function siteOrigin(env = process.env) {
  return env.DOCS_SITE || (env.GITHUB_REPOSITORY_OWNER
    ? `https://${env.GITHUB_REPOSITORY_OWNER.toLowerCase()}.github.io`
    : 'https://molik-labs.github.io');
}
