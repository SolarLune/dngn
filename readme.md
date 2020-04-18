
# dngn

![dngn_v0.2](https://user-images.githubusercontent.com/4733521/48660612-68e94480-ea19-11e8-8f4d-b378fa64dabe.gif)

[GoDocs](https://pkg.go.dev/github.com/SolarLune/dngn?tab=doc)

## What is dngn?

dngn is a golang library specifically created to help make generating random maps easier.

## Why is it called that?

Because it's short, simple, and to the point. Also, vwls r vrrtd.

## Why did you create dngn?

Because I needed to do random map generation for a game and didn't seem to find a library around anywhere.

And so, here we are.

## How do I install it?

Just go get it and import it in your game application.

`go get github.com/SolarLune/dngn`

Or import and use it directly in your code with a `go.mod` file in your project directory, and it will go ahead and automatically download it and use it for your project. [See the Go wiki.](https://github.com/golang/go/wiki/Modules#quick-start)

## How do I use it?

dngn is based around Rooms and Selections. A Room contains a rune array, representing the Room's data. You can either manipulate the data array manually, or use Selections to grab a selection of cells in the Room to alter. 

To start off with using dngn, you can just create a Room, and then use one of the included `Generate` functions to generate the data:

```go
import "github.com/SolarLune/dngn"

var GameMap *dngn.Room

func Init() {

    // This line creates a new Room. The size is 10x10.
    GameMap = dngn.NewRoom(10, 10)

    // This will select the cells the map has, and then fill the selection with "x"s.
    GameMap.Select().Fill('x')

    // Selections are structs, so we can store Selections in variables to store the "view" of the data.
    selection := GameMap.Select()

    // This will run a drunk-walk generation algorithm on the Room. It starts at a random point
    // in the Room, and walks around the Room, placing the value specified (0, in this case)
    // until the room is the percentage provided (0.5, or 50%, in this case) filled.
    GameMap.GenerateDrunkWalk(' ', 0.5)

    // This function will degrade the map slightly, making cells with a ' ' in them randomly turn into a cell with a 'x',
    // depending on how heavily the cell is surrounded by 'x's.
    GameMap.Select().ByRune(' ').Degrade('x')

    // Room.DataToString() will present the data in a nice, easy-to-understand visual format, useful when debugging.
    fmt.Println(GameMap.DataToString())

    // Now we're done! We can use the Room.

}

```

Selections can also be powerful, as they allow you to easily select cells to manipulate. You can also chain Selection filtering functions together. As an example, say you wanted to randomly change a small percentage of floor tiles (' ') into trap tiles ('z'). You could easily do this with Selections, like so:

```go
    GameMap.Select().ByValue(' ').ByPercentage(0.1).Fill('z')
```

---

And that's about it! There are also some nice additional features to make it easier to handle working with and altering Rooms.

## Wait... How do I actually LOOK at it?

Welp.

So dngn just does map generation - it doesn't handle visualization / rendering of the map. For that, you can use another framework, like [pixel](https://github.com/faiface/pixel), [Ebiten](https://github.com/hajimehoshi/ebiten), [raylib-goplus](https://github.com/Lachee/raylib-goplus), or [go-sdl2](https://github.com/veandco/go-sdl2).

Soooo that's about it. If you want to see more info or examples, feel free to examine the main.go and world#.go tests to see how a couple of quick example tests are set up.

[You can check out the GoDoc link here, as well.](https://pkg.go.dev/github.com/SolarLune/dngn?tab=doc)

You can also run the example by simply running the example from the project's root directory:

```
$ go run ./example/
```

## Dependencies?

For the actual package, there are no external dependencies. dngn just uses the built-in "fmt" and "math" packages.

For the tests, dngn requires [Ebiten](https://github.com/hajimehoshi/ebiten) to create the window, handle input, and draw the shapes.
