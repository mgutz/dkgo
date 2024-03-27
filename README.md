# dkgo

Utility commands for [dk tiling WM](https://bitbucket.org/natemaia/dk)

To Install

    go install github.com/mgutz/dkgo

## Commands

| command     | description                                                                                                                                                                                            |
| ----------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| cycle-ws    | Cycle through dynamic workspaces. MacOS behavior.                                                                                                                                                      |
| rofi-master | Swaps a window to master by selecting from window list. Requires [rofi](https://github.com/davatorium/rofi).                                                                                           |
| swap-master | Smarter swap. Swap by tile index, window ID. If empty then it swaps with last or current focused. TIP: Execute swap-master twice to quickly switch between master and last. Works on floating windows. |
| status      | Pretty prints dkcmd status.                                                                                                                                                                            |

## License

MIT, see [LICENSE](LICENSE)
