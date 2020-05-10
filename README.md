# wp - wallpaper manager

wp is a command-line wallpaper manager for Linux.
It is mostly suitable for minimal window managers (i3, dwm, bspwm, etc) that do not manage wallpapers themselves 
(and, more importantly, do not override the background of the X root window).

# Features

* JPEG, PNG, and GIF file formats.
* Random wallpapers from [Unsplash](https://unsplash.com).
* Multiple displays are supported via xinerama.
* Automatic wallpaper change with cron.
* Self-contained, does not require external tools.

# Installation

```bash
GO111MODULE=on go get github.com/mkuznets/wp/cmd/wp
```

# Configuration

`wp` expects a YAML config file in the following locations, in order:

* path from `-c` or `--config`
* `$XDG_CONFIG_HOME/wp.yaml`
* `$HOME/.config/wp.yaml`

```sh
# Create a config file with default values
wp config --defaults > $XDG_CONFIG_HOME/wp.yaml
```


# Usage

```shell
# Set wallpaper from a local file
wp file /path/to/image.jpg

# Set a random wallpaper from Unsplash (if configured)
wp unsplash

# Set new wallpaper from Unsplash if the current one is older than `tick.ttl`
# (Supposed to be run regularly via cron.)
wp tick

# Save the current wallpaper into a configured directory (`fs.gallery`)
wp save

# Restore the most recent wallpaper
# (Supposed to be run from xinit/xprofile scripts.)
wp init
```
