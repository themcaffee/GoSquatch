[_metadata_:title]:- "Github Action"
[_metadata_:layout]:- "index"

# Github Action

GoSquatch provides a Github Action on the Marketplace. You can check it out [here!](https://github.com/marketplace/actions/gosquatch)

## Inputs

#### `srcDir`

The source directory to pull the markdown and templates from. Default `"src"`.

## Configuration

GoSquatch is configured with a `.squatch` in the folder `srcDir`. This file is not required and the action will just fine without it. However,
if you need additional configuration options then it is available.

## Example Usage

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
        uses: themcaffee/GoSquatch@1.0.28-beta
        with:
          srcDir: 'src'
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