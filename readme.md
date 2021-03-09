
# dngn

![dngn_v0.2](https://user-images.githubusercontent.com/4733521/48660612-68e94480-ea19-11e8-8f4d-b378fa64dabe.gif)

[pkg.go.dev docs](https://pkg.go.dev/github.com/SolarLune/dngn?tab=doc)

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

dngn is based around Layouts and Selections. A Layout contains a rune array, which is the internal data. You can either manipulate the data array manually, or use Selections to grab a selection of cells in the Layout to alter. Layouts have Generate functions to generate a game map; they each take different arguments to influence how the function generates a map. To start off, you can just create a Layout, and then use one of the included `Generate` functions to generate the data:

```go
import "github.com/solarlune/dngn"

var GameMap *dngn.Layout

func Init() {

    // This line creates a new *dngn.Layout. The size is 10x10.
    GameMap = dngn.NewLayout(10, 10)

    // This will generate a map using BSP map generation and some sensible default settings.
    GameMap.GenerateBSP(dngn.NewDefaultBSPOptions())

    // Layout.DataToString() will present the data in a nice, easy-to-understand visual format, useful when debugging.
    fmt.Println(GameMap.DataToString())

    // Now we're done! We can visualize and use the Layout.

}

```

As seen above, Selections can be chained together. As an example, say you wanted to randomly change a small percentage of floor tiles (' ') into trap tiles ('z'). You could easily do this with Selections, like so:

```go
    GameMap.Select().FilterByValue(' ').FilterByPercentage(0.1).Fill('z')
```

---

And that's about it! There are also some nice additional features to make it easier to handle working with and altering Layouts.

## Wait... How do I actually LOOK at it?

dngn just does map generation - it doesn't handle visualization / rendering of the map. For that, you can use another framework, like [pixel](https://github.com/faiface/pixel), [Ebiten](https://github.com/hajimehoshi/ebiten), [raylib-goplus](https://github.com/Lachee/raylib-goplus), or [go-sdl2](https://github.com/veandco/go-sdl2).

That's about it. You can run the example by simply running the example from the project's root directory:

```
$ go run ./example/
```

[pkg.go.dev docs](https://pkg.go.dev/github.com/SolarLune/dngn?tab=doc)

## Dependencies?

For the actual package, there are no external dependencies. dngn just uses the built-in "fmt" and "math" packages.

For the tests, dngn requires [Ebiten](https://github.com/hajimehoshi/ebiten) to create the window, handle input, and draw the shapes.
