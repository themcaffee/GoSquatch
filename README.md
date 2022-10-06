# GoSquatch

A super fast Github Action that converts markdown into a static HTML site. This is super useful for personal blogs and project documentation
to keep pages in standard markdown but also be able to host through Github Pages (or other hosting providers). GoSquatch uses [native golang templating](https://pkg.go.dev/text/template) and [gomarkdown/markdown](https://github.com/gomarkdown/markdown) to handle markdown parsing.

_How fast?_ 
GoSquatch takes about 3 seconds on Github Actions to execute. This allows with checking out the code and publishing it to Github Pages to only take around 
20 - 30 seconds in total execution time. Check out this repo's [Actions](https://github.com/themcaffee/GoSquatch/actions) for real examples of performance.


_Why is it so fast?_ 
First, the docker container for this action is a [seperate repository](https://github.com/themcaffee/GoSquatchDocker) that builds an extremely lean docker 
image with only an alpine image and a small binary program file. This allows for this action to pull a very small image hosted on Github that only takes 
about 3s to pull down. Second, because this is written in Go this allows for a tight binary with super fast execution with the minimal depencies built in 
the binary. The step to build the pages varies depending 
on size but will generally be less than 1 second. 

[Check out our documentation built with GoSquatch!](https://mitchmcaffee.com/GoSquatch/)

## Inputs

#### `srcDir`

The source directory to pull the markdown and templates from. Default `"src"`. Not required.

#### `distDir`

The distribution directory where the static html files will be put. Default `"dist"`. Not required.

#### `ignoreFolders`

Comma seperated list of folders to ignore. Any folders starting with "." will always be ignored. Default `"node_modules"`. Not required.

#### `ignoreFiles`

Comma seperated list of files to ignore. Default `"LICENSE, yarn.lock"`. Not required.

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

## Note on text/template usage and Safety

This project uses text/template over html/template because GoSquatch is generating HTML over just formatted text placed in an HTML document. Do not use 
this on untrusted sources because any .js files will be copied in to the folder that gets uploaded to Github Pages. This is by design so that javascript 
sources can be included in the website. If you want to ignore files or folders, there are inputs to the action that can be used.


## License

The scripts and documentation in this project are released under the [MIT License](https://github.com/themcaffee/GoSquatch/blob/main/LICENSE).