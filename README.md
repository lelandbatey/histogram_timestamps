# What's the idea?

basic idea is to also write a Go program which will do the actual parsing and
binning of the time-series data, then render a single static HTML page with two
script tags, one with the binned data in JSON form assigned to a global
variable, the second containing the full contents of `bundle.js`. They should
be on that page in that order so that we can be sure about the order of
interpretation. The code in the `bundle.js` will always look for data in the
global variable that'll be defined in the first script tag.

# How to rebuild bundle.js?

bundle.js is built from main.js . To build `bundle.js`, use the following command:
```
$(npm bin)/browserify main.js -o bundle.js
```

- I used this linked document as basis for how to structure this project, since I want it to be possible to "self-contain".  https://medium.com/jeremy-keeshin/hello-world-for-javascript-with-npm-modules-in-the-browser-6020f82d1072

