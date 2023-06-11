<!--Flags settup of the application> -->

# Image to ASCII

This is a simple application that converts an image to ASCII art.

### supported formats

- gif

future support

- jpg
- png
- bmp


## How to use

```bash
$ go run main.go -w=25 -h=25 -i=localfile.gif
```
```bash
$ go run main.go -w=75 -h=100 -i=https://upload.wikimedia.org/wikipedia/commons/5/5a/Rotating_Tux.gif 
```

## how to build

```bash
$ go build -o image2ascii main.go
```

## Flags

| Flag | Description | Default | Needed |
| --- | --- | --- | --- |
| `-w` | width of the image | 25 | No |
| `-h` | height of the image | 25 | No |
| `-i` | input source image , it can be a localfile or external url | empty | Yes |


## Recomended terminals

- Kitty - this terminal is the best in the test for the moment 
- KDE Konsole 