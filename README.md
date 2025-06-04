# GoSquatch

GoSquatch is a fast GitHub Action that converts Markdown into a static HTML site.
It is useful for personal blogs and project documentation, letting you keep pages in standard Markdown while hosting them through GitHub Pages or other providers.
GoSquatch uses [native Go templating](https://pkg.go.dev/text/template) and [gomarkdown/markdown](https://github.com/gomarkdown/markdown) to handle Markdown parsing.
It includes a live server so you can preview your site locally before publishing.
See the [performance documentation](https://mitchmcaffee.com/GoSquatch/performance) for details on execution speed.

[Check out our documentation built with GoSquatch!](https://themcaffee.github.io/GoSquatch/)


## Getting Started

### Install the CLI

Run the following commands to install GoSquatch using `apt`:

```bash
curl -s --compressed "https://www.mitchmcaffee.com/GoSquatch/ppa/KEY.gpg" | gpg --dearmor | sudo tee /etc/apt/trusted.gpg.d/gosquatch.gpg > /dev/null
sudo curl -s --compressed -o /etc/apt/sources.list.d/gosquatch_list_file.list "https://www.mitchmcaffee.com/GoSquatch/ppa/gosquatch_list_file.list"
sudo apt update
sudo apt install gosquatch
```

To preview your site locally run:

```bash
gosquatch -live-server -src-dir=./
```

### GitHub Action

Add GoSquatch to your workflow to build pages during CI:

```yaml
- uses: actions/checkout@v2
- name: Build pages
  uses: themcaffee/GoSquatch@1.0.28-beta
```

The action reads Markdown from `src` and outputs static files to `dist` ready for hosting.


## License

The scripts and documentation in this project are released under the [MIT License](https://github.com/themcaffee/GoSquatch/blob/main/LICENSE).
