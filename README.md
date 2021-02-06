# tcon

A simple wrapper for tcell-based text-only interfaces.

## Requirements

-   `go get github.com/gdamore/tcell`

## What and why?

I wanted a terminal-based environment with a fixed-position command line and a large text area to display text. Other packages were mimicking GUIs in text mode, but that's more than I needed for some of my projects.

## Features

-   Command box with editable input and history.
-   Text output with word wrap and less-mode.
-   Title and status lines.
-   Callbacks for input and command handling.

## Keyboard shortcuts

There is a separate handler for different types of input. These are:

-   OnCommandFunc: This callback receives a string with the command buffer. The rest is up to the user.
-   OnRuneFunc: Filter callback to allow special handling while typing. Optionally return 0 to skip symbols.
-   OnFunc: The generic callback, used with SetTabFunc() and SetEscFunc().
-   OnCtrlFunc: Handler for control keys.

NOTE: Remember to set up a key or command to call the Quit() method.

## TODO

(Possible) future enhancements:

-   Coloured status and mode
-   Command line manipulation in the filter?
-   Built-in TAB-completion?
-   Figure out why resizing sometimes messes up things
