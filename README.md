# gosass #

## About ##

Gosass is a fast sass compiler. We made this because the normal ruby compiler
ran too slowly, and the native compiler (`sassc`) didn't have all the features
we needed (directory compilation and watching.)

## Installation ##

Gosass depends on `sassc`, so
[install that first](https://github.com/sass/sassc). On mac, you can use
homebrew: `brew install sassc`.

Then install `gosass`:

```bash
go get github.com/dailymuse/gosass
pushd $GOPATH/src/github.com/dailymuse/gosass
make install
popd
```

## Installation in Docker ##

### Compile gosass

```bash
sudo docker run --rm -v "$PWD":/usr/src/myapp -w /usr/src/myapp golang:1.4 bash -c "go get ./... ; CGO_ENABLED=0 go build -a -installsuffix cgo && mv myapp gosass"
```

### Build Docker image
```bash
sudo docker build -t local/gosass .
```

### Run Docker container
```bash
sudo docker run -ti --rm -v $PWD:$PWD -w $PWD local/gosass -input webapp/scss/ -output webapp/static/css/ -style compressed
```

Append `-watch` if wanted.

## Example Usage ##

```bash
# Compile an individual file
gosass -input file.scss -output file.css

# Compile a directory
gosass -input sass/ -output css/

# Live-compile a directory
gosass -input sass/ -output css/ -watch
```
