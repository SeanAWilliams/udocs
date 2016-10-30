# TODO

> no promises

### Issues

- search page rendering should be ajax
- there is a padding discrepancy between document and inner rendering
- make default template dir overridable
- anchor tag goto is about 50 px's off
- handler tests... all of them
- Encapsulate mutatation of settings and instead use env vars, for dynamic restaging, and then just log when you do it (for transparency)
- Make error messages less stack-based, and more human-readable
- mobile view collpases sidebar and search out of view

### Features

- update to newer version of treemux, and utilize Context features
- implement default index.html, and support for generating one
- SUMMARY.md should support absolute paths, and URL paths, to markdown files
- Add a single-repo view for the sidebar that will default expand the sidebar if only a single docs directory is being hosted
- Add a README gif showing the CLI in-action
- Implement '-v, --verbose' flags that enable log output. By default, it should be off (which it currently is not)