# Configuration

## File based configuration

In the root of your repository, add a new file `.gosquatch` to define GoSquatch's configuration. The available options are:

- `dist`: Directory to output built files. This folder will be created if it does not exist.
- `IgnoreFiles`: List of comma separated file names to ignore when building. These files will not be copied over into the output directory.
- `IgnoreFolders`: List of comma separated folder names to ignore when build. These folders and their contents will not be copied over into the output directory.

Example `.gosquatch` file:

```yaml
dist: dist
ignoreFolders: node_modules,static/tmp
ignoreFiles: README.md
```
