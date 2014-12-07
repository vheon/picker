# Picker

A terminal only fuzzy finder. It is inspired by [Selecta][] and the others but written in [golang][]

## Rationale

I really like the way [Selecta][] works, it do not depend on anything other than an ANSI terminal and do not use the entire screen so I still have some context.
So why write yet another fuzzy finder where [Selecta][] is already there beside [gof][], [fzf][]?
Well:
* [Selecta][] is a little too slow for my taste when dealing with very large input due to its internal design
* [gof][] it print on the entire screen, leaving me without context and is also slow with very large input maybe due to the fact that it use regexp for matching and printing the matching part of the screen
* [fzf][] use the entire screen for printing and have more options and features than I need

Plus and more important, I wanted to write one, and used this tool as an exercise for learning [golang][].

[golang]: http://golang.org
[Selecta]: https://github.com/garyberngardt/selecta
[gof]: https://github.com/mattn/gof
[fzf]: https://github.com/junegunn/fzf

## Usage

Picker as most Unix tool takes its input from stdin and write the selection to stdout:

```
$ vi $(find '.' | picker)
```

* `<Enter>` key select the candidate
* `<Esc>` or `<C-c>` cancel the operation
* `<C-n>` select down
* `<C-p>` select up
