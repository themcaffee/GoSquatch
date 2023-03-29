# Configuration

## File based configuration

In the root of your repository, add a new file `.gosquatch` to define GoSquatch's configuration. The available options are:

- `dist`: Directory to output built files. This folder will be created if it does not exist.
- `IgnoreFiles`: List of comma seperated file names to ignore when building. These files will not be copied over into the output directory.
- `IgnoreDirectory`: List of comma seperated folder names to ignore when build. These folders and their contents will not be copied over into the output directory.