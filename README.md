# Squatch

A super fast Github Action that converts markdown into a static HTML site.

## Inputs

#### `srcDir`

The source directory to pull the markdown and templates from. Default `"src"`.

#### `distDir`

The distribution directory where the static html files will be put. Default `"dist"`.

### Example Usage

```
uses: themcaffee/GoSquatch@v1
with:
    srcDir: 'src'
    distDir: 'dist'
```