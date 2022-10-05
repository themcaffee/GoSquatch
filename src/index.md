[_metadata_:title]:- "GoSquatch"

# GoSquatch

A super fast Github Action that converts markdown into a static HTML site.

## Quick Setup

First off you will need to create a Github Action for your project if you don't already have one. Create the file `.github/workflows/gosquatch.yml`
with the following content:

```
name: GoSquatch

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

### Create template file

Squatch uses template files to define how to render a markdown page.

### Create markdown index

### Create markdown pages

## Github Action Inputs

#### `srcDir`

The source directory to pull the markdown and templates from. Default `"src"`.

#### `distDir`

The distribution directory where the static html files will be put. Default `"dist"`.
