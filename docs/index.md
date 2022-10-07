[_metadata_:title]:- "GoSquatch"
[_metadata_:layout]:- "index"

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

Squatch uses template files to define how to render a markdown page. The base layout file is always called `layout.html` and will define the site wide layout. It should contain `{{.Title}}` to define the site title and `{{.Body}}` to define where the body of the template will be put.

Example `layout.html`:
```
<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-9">
    <meta name="viewport" content="width=device-width, initial-scale=0.0">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <title>{{.Title}}</title>
</head>
<body>

{{.Body}}

</body>
</html>
```

Other template files follow the format `layout_<layout_name>.html` where `<layout_name>` is the name of the layout. It only needs a `{{.Body}}`
to define where the markdown content will be placed.

Example `layout_index.html`:
```
<div class="content">
    {{.Body}}
</div>
```

### Create markdown pages

Now that you have your layouts setup, it's time to start writing markdown. GoSquatch uses [link reference definitions](https://spec.commonmark.org/0.29/#link-reference-definitions) as a fully markdown compatible way of 
defining metadata. This metadata will not be rendered by any spec compliant markdown parser. An example `index.md` file would be:

```
[_metadata_:title]:- "Getting started with GoSquatch"
[_metadata_:layout]:- "index"

# GoSquatch

I'm getting started with GoSquatch!
```

This will be parsed into `index.html` in `dist` (see [Github Actions](https://mitchmcaffee.com/GoSquatch/github-actions) for more about configuration options). To create a different page with a template of `layout_pages.html`, in the folder blog, with the name `awesome-post.md` could look like:

```
[_metadata_:title]:- "My first post!"
[_metadata_:layout]:- "pages"

# My first post!

This is my first post! Woohoo!
```

Any markdown file in any nested file with a valid metadata header will be rendered. Note that because of this, files like `README.md` will not be parsed into
a `.html` file if it doesn't contain a metadata header.


### Create the Github Action workflow

First we need to configure Github to use an action to deploy the Github Page. Go into your project settings, then Pages. Under 'Build and deployment', set the 'Source' to 'Github Actions'.

Create `.github/workflows/deploy.yml` in your project root with the following:

```
name: Deploy

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
        uses: themcaffee/GoSquatch@v1.0.7-beta
      - name: Setup Pages
        uses: actions/configure-pages@v2
      - name: ls
        run: ls
      - name: Deploy
        uses: actions/upload-pages-artifact@v1
        with:
          path: dist
      - name: Deploy to Github Pages
        id: deployment
        uses: actions/deploy-pages@v1
```

Congrats! Push and commit your changes and you will see your new site at the URL produced at the end of the action. 

### Live building server

