# Squatch

A super fast Github Action that converts markdown into a static HTML site.

## Inputs

#### `src-dir`

The source directory to pull the markdown and templates from. Default `"src"`.

#### `dist-dir`

The distribution directory where the static html files will be put. Default `"dist"`.

### Example Usage

```
uses: actions/Squatch@v1
with:
    srcDir: 'src'
    distDir: 'dist'
```