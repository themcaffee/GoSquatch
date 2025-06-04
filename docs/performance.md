[_metadata_:title]:- "Performance"
[_metadata_:layout]:- "index"

# Performance

GoSquatch typically completes its action in about 3 seconds on GitHub Actions. Including checkout and publishing to GitHub Pages, an entire workflow often finishes in 20–30 seconds. See the [GoSquatch-template actions](https://github.com/themcaffee/GoSquatch-template/actions) for real-world examples.

## Why it is fast

The action uses a slim Docker image built from Alpine with a small Go binary. Pulling the image only takes a few seconds and Go's statically compiled binary executes quickly with minimal dependencies. Building pages generally takes less than a second.
