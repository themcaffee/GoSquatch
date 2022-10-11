# GoSquatch

A super fast Github Action that converts markdown into a static HTML site. This is super useful for personal blogs and project documentation
to keep pages in standard markdown but also be able to host through Github Pages (or other hosting providers). GoSquatch uses [native golang templating](https://pkg.go.dev/text/template) and [gomarkdown/markdown](https://github.com/gomarkdown/markdown) to handle markdown parsing. It includes a live building
server so you can easily see your site locally before publishing it.

_How fast?_ 

GoSquatch takes about 3 seconds on Github Actions to execute. This allows with checking out the code and publishing it to Github Pages to only take around 20 - 30 seconds in total execution time. Check out the [GoSquatch-template's Github Actions](https://github.com/themcaffee/GoSquatch-template/actions) for real examples of performance. 


_Why is it so fast?_ 

First, there is a seperately built and published docker image `themcaffee/gosquatch` that builds an extremely lean docker 
image with only an alpine image and a small binary program file. This allows for this action to pull a very small image hosted on Github that only takes 
about 3s to pull down. Second, because this is written in Go this allows for a tight binary with super fast execution with the minimal depencies built in 
the binary. The step to build the pages varies depending on size but will generally be less than 1 second. 

[Check out our documentation built with GoSquatch!](https://mitchmcaffee.com/GoSquatch/)


## License

The scripts and documentation in this project are released under the [MIT License](https://github.com/themcaffee/GoSquatch/blob/main/LICENSE).
