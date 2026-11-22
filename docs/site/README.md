# FCC documentation site

This directory contains the static source for the FCC guide. The site covers
installation, providers, clients, APIs, messaging, security, maintenance and
contribution in English and Vietnamese. It links to the project README for the
full provider catalog and platform-specific details; update both translations
and the README when behavior changes.

## Build and preview

From the repository root, with Python 3.9 or newer:

```sh
python3 -m unittest tests.scripts.test_docs_site -v
python3 scripts/build_docs_site.py
python3 -m http.server 8000 --directory site-build
```

Open <http://localhost:8000/> on the machine running the server. To avoid
accidentally overwriting files, the builder requires a new or empty output
directory. Remove the generated, ignored `site-build/` directory before a
subsequent build, or use `--output` with another empty directory. It copies the
site source and selected existing screenshots, creates `.nojekyll`, and checks
local links, fragments, duplicate IDs and image alt attributes. It does not
verify external links, rendered appearance or product behavior. No additional
packages are needed.

## Deploy

The `Documentation site` workflow checks every pull request. Pushes to `main`
and manually dispatched runs on `main` build a Pages artifact and deploy it
using GitHub Actions. A repository administrator must select **Settings → Pages →
Build and deployment → Source: GitHub Actions**, and any required environment
approval must be granted. See the workflow's deploy job for the actual page
URL. This workflow never publishes a pull request or working branch, and adding
it to a branch does not itself publish a site or change repository settings.

Links to project source and reference material point to the `main` branch of
`MoLik-Labs/FCC`; they represent the published documentation after these
changes reach that branch. The installation commands explicitly refer to the
upstream project and are not an installer for this fork.
