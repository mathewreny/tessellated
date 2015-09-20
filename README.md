# Tessellated
This is a fun little program I wrote that creates a randomized tessellated svg image.
It's written in go and is used on my website http://reny.io

# Package
To use tessellated as a library, simply call the Triangle function and pass in an io.Writer.  
In the future I plan on creating a rectangle and Hexagon method.

    tessellated.Triangle(os.Stdout, 1920, 1200)


# Command Line
To generate an svg using the command line, provide a width and height.  The default width and height are 1000px. 
`0 < width` and `0 < height`

    $ ./tessellated --width 1920 --height 1200 > image.svg


**Toy Server**: To run a server that serves randomly generated triangle svgs enter the following into the command line

    $ ./tessellated --http 12345

then in your browser visit

    http://localhost:12345/triangle.svg?width=1920&height=1200

where width and height are integers `0 < width <= 5120` and `0 < height <= 2880`.
The width and height must be provided
