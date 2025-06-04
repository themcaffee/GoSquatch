# GoSquatch

GoSquatch is a fast Github Action that converts markdown into a static HTML site. This is useful for personal blogs and project documentation to keep pages in standard markdown while hosting them through Github Pages (or other providers). GoSquatch uses [native golang templating](https://pkg.go.dev/text/template) and [gomarkdown/markdown](https://github.com/gomarkdown/markdown) to handle markdown parsing. It includes a live building server so you can easily preview your site locally before publishing it. See the [performance documentation](https://mitchmcaffee.com/GoSquatch/performance) for details on execution speed.

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

The action reads markdown from `src` and outputs static files to `dist` ready for hosting.


## License

The scripts and documentation in this project are released under the [MIT License](https://github.com/themcaffee/GoSquatch/blob/main/LICENSE).
