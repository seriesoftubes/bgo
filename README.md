# Bgo: Backgammon engine written in Go

### Quick start
 - Clone this repo into $GOPATH/src/github.com/seriesoftubes: 
 ```sh
 mkdir -p $GOPATH/src/github.com/seriesoftubes && cd $GOPATH/src/github.com/seriesoftubes && git clone https://github.com/seriesoftubes/bgo.git && cd bgo
 ```
 - Build main.go
```sh
go build main.go
```
- Play against untrained AI opponent by entering moves like `X;a1;m5` (the "X" is your player name, "a1" means move X's checker on the "a" slot by 1, "m5" means move X's checker on the "m" slot by 5)
```sh
./main -skip_training
```

### Training the AI opponent
This can be done by adjusting the training parameters via command line flags and interactively adjusting settings at runtime.
- Run with reasonable flags:
```sh
./main
```
- Further train a pre-trained opponent:
```sh
./main -epsilon=0.3 -config_infile='~/Desktop/bgo/bgo_nnet.json' -config_outfile='~/Desktop/ai/agent2.json'
```