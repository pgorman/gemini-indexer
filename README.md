Gemini Indexer
========================================

Gemini Indexer is a simple utility to generate an index page for a directory of [Gemini](https://gemini.circumlunar.space/) files.

The workflow is:

1. Create or upload your Gemini files to the directory where your Gemini server makes them available (e.g., `~/public_gemini/`).
2. Run `gemini-indexer` (either manually or as a periodic cron job) to create an `index.gmi` file listing the contents of that directory.

These are some of the important command-line options for `gemini-indexer`.
Run `gemini-indexer --help` for a complete list.

```
--dotfiles    include dotfiles, like '..' and '.git' in the index
--ignore      comma-separated list of files to not include in the index
--indir       path to the directory of files to index (default: current directory
--outfile     where to write the index file (default: stdout)
--template    template file for the index (see https://pkg.go.dev/text/template)
--title       title displayed on index page, like "Jane Smith's Gemlog"
```

You are strongly encouraged to set `--title` to a name unique to your site.

Gemini Indexer will attempt to set a date for any `*.gmi` or `*.gemini` files in order to produce an index compatible with the [Gemini subscription specification](gemini://gemini.circumlunar.space/docs/companion/subscription.gmi).
If the file names begin or end with a YYYY-MM-DD date, like `2020-01-02-hello.gmi` or `hello-2020-01-02.gmi`, that will be used as the file date in the link text; if no such date is include in the file name, the modification time of the file will be used.

The remaining file name, minus the date, will be used as the label of the link, with any underscores (_) replaced by spaces.

For example, a file named `2020-12-21-My_Tropical_Vacation.gmi` results in this link line:

```
=> 2020-12-21-My_Tropical_Vacation.gmi 2020-12-21 My Tropical Vacation
```

If you want to customize the output of Gemini Indexer, see `example.tmpl` as the basis for your custom template.


Copyright
----------------------------------------

Gemini Indexer copyright 2020 Paul Gorman.

Licensed under the GPL. See LICENSE.txt for details.