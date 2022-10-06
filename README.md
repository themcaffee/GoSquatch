# GoSquatch

A super fast Github Action that converts markdown into a static HTML site. This is super useful for personal blogs and project documentation
to keep pages in standard markdown but also be able to host through Github Pages (or other hosting providers).

[Check out our documentation!](https://mitchmcaffee.com/GoSquatch/)

## Inputs

#### `srcDir`

The source directory to pull the markdown and templates from. Default `"src"`.

#### `distDir`

The distribution directory where the static html files will be put. Default `"dist"`.

### Example Usage

```
name: Docs

on:
  push:
    branches:
      - main

permissions:
  contents: read
  pages: write
  id-token: write

concurrency:
  group: "deploy"
  cancel-in-progress: true

jobs:
  deploy:
    runs-on: ubuntu-latest
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    steps:
      - uses: actions/checkout@v2
      - name: Build pages
        uses: themcaffee/GoSquatch@v1-beta
        with:
          srcDir: 'src'
          distDir: 'docs'
      - name: Setup Pages
        uses: actions/configure-pages@v2
      - name: ls
        run: ls
      - name: Deploy
        uses: actions/upload-pages-artifact@v1
        with:
          path: docs
      - name: Deploy to Github Pages
        id: deployment
        uses: actions/deploy-pages@v1
```