
# dngn

![dngn](https://user-images.githubusercontent.com/4733521/47975754-fadd6e80-e063-11e8-932c-f05d8110eddf.gif)

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

## How do I use it?

dngn is based around Rooms. A Room contains an int array, representing the Room's data. You can either manipulate the data array manually, or use the included functions to alter the data. Rooms can also contain other Rooms, forming a hierarchy - in this case, a child Room will point to the root Room's map data array. 

To start off with using dngn, you can just create a Room, and then use one of the included Generate functions to generate the data:

```go

var GameMap *dngn.Room

func Init() {

    // This line creates a new Room. The position is 0, 0, and is ignored if the Room is a 
    // root (a root Room doesn't have a parent). The size is 10x10, and the parent is nil (so it is a root).
    GameMap = dngn.NewRoom(0, 0, 10, 10, nil)

    // This will fill the map with "1"s.
    GameMap.Fill(1)

    // This will run a drunk-walk generation algorithm on the Room. It starts at a random point 
    // in the Room, and walks around the Room, placing the value specified (0, in this case) 
    // until the room is the percentage provided (0.5, or 50%, in this case) filled.
    GameMap.GenerateDrunkWalk(0, 0.5)

    // This function will degrade the map slightly, making cells with a 0 in them randomly turn into a cell with a 1 in it. 
    // If the cell is next to the target value (1), then it's more likely to turn into a 1. 
    // If it isn't surrounded on any sides by a 1, then it won't.
    GameMap.Degrade(0, 1)

    // Now we're done! We can use the Room.

}

```

And that's about it! There are also some nice additional features to make it easier to handle working with and altering Rooms.

## Wait... How do I actually LOOK at it?

Welp. 

So dngn just does map generation - it doesn't handle visualization / rendering of the map. For that, you can use another framework, like [pixel](https://github.com/faiface/pixel), [Ebiten](https://github.com/hajimehoshi/ebiten), [raylib](https://github.com/gen2brain/raylib-go), or [go-sdl2](https://github.com/veandco/go-sdl2).

Soooo that's about it. If you want to see more info or examples, feel free to examine the main.go and world#.go tests to see how a couple of quick example tests are set up.

[You can check out the GoDoc link here, as well.](https://godoc.org/github.com/SolarLune/dngn/dngn)

## Dependencies?

For the actual package, there are no external dependencies. dngn just uses the built-in "fmt" and "math" packages.

For the tests, dngn requires veandco's sdl2 port to create the window, handle input, and draw the shapes.

## Shout-out Time!

Props to whoever made arcadepi.ttf! It's a nice font.

Thanks a lot to the SDL2 team for development.

Thanks to veandco for maintaining the Golang SDL2 port, as well!
