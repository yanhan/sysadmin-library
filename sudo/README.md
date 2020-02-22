# About

Simple examples of lines you can add to `/etc/sudoers` or `/etc/sudoers.d` to grant some groups / users full / limited access.

Code is in the [examples](examples) file.


## Protip

Please use `visudo` to test out individual commands, especially before adding files to `/etc/sudoers.d`.



## Change visudo editor

```
sudo update-alternatives --config editor
```

then select `vim.basic` or `vim.tiny`.
