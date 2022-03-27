# Checkers

 ![Checkers](/doc/animation.gif)

Simple checkers, playing arround with minimax heuristics.

It can be launched in full ai vs ai mode with the -ai flag.
The main goal here wasnt the gameitself but rather exploring minimax and heuristics.

I ended up using adapted heurstics from https://github.com/kevingregor/Checkers/blob/master/Final%20Project%20Report.pdf for now,
but still exploring possibilites.

Due to the protection heuristic some AI vs AI games will result in some wall hugging.

If you might wonder why each ai run yields different results despite no shuffelinig, 
its because the go map is intentionally random when iterated with range,
thus if there are many equal moves it will pick a random one and there is no need to shuffle.

## Building

```
go mod download
go build main.go
```


## Used Packages

https://github.com/faiface/pixel  - used to draw the Board

## License

BSD-2-Clause