# TODO

> no promises

### Issues

- compile all static files into the binary, to clean up build/package (https://github.com/jteeuwen/go-bindata)
- all HTML pages should be entirely pre-rendered, so that they serve faster. currently, they are still being templated into the response
- anchor tag goto is about 50 px's off
- running udocs serve -d <dir>, where dir is not docs dir, gives weird error message
- bin/install.sh is needlessly specific, and often redudant in regards to Dockerfile
- handler tests... all of them
- Encapsulate mutatation of settings and instead use env vars, for dynamic restaging, and then just log when you do it (for transparency)
- Disable file server currently allowing static files directory to be hit from browser
- Make error messages less stack-based, and more human-readable

### Features

- update to newer version of treemux, and utilize Context features
- implement default index.html, and support for generating one
- SUMMARY.md should support absolute paths, and URL paths, to markdown files
- Make the primary color (green, currently) configurable
- Add a single-repo view for the sidebar that will default expand the sidebar if only a single docs directory is being hosted
- Add a README gif showing the CLI in-action
- Make the Git SSH key configurable
- Implement '-v, --verbose' flags that enable log output. By default, it should be off (which it currently is not)