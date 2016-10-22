# TODO

> no promises

### Issues

- make default template dir overridable
- anchor tag goto is about 50 px's off
- running udocs serve -d <dir>, where dir is not docs dir, gives weird error message
- handler tests... all of them
- Encapsulate mutatation of settings and instead use env vars, for dynamic restaging, and then just log when you do it (for transparency)
- Make error messages less stack-based, and more human-readable

### Features

- update to newer version of treemux, and utilize Context features
- implement default index.html, and support for generating one
- SUMMARY.md should support absolute paths, and URL paths, to markdown files
- Add a single-repo view for the sidebar that will default expand the sidebar if only a single docs directory is being hosted
- Add a README gif showing the CLI in-action
- Implement '-v, --verbose' flags that enable log output. By default, it should be off (which it currently is not)