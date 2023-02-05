package main

import (
    "fmt"
    "image"
    "image/color"
    "math"
    "image/png"
    "os"
)

func PowInt(base, exponent int) int {
    out := 1
    for i := 0; i < exponent; i++ {
        out *= base
    }
    return out
}

type Sample []bool

func getSamples(k int) []Sample {
    n := PowInt(2, k)
    samples := make([]Sample, n)
    for i := 0; i < n; i++ {
        samples[i] = make([]bool, k)
    }
    truth := false
    flip := n/2
    for i := 0; i < k; i++ {
        for j := 0; j < n; j++ {
            if j%flip == 0 {
                if truth {
                    truth = false
                } else {
                    truth = true
                }
            }
            samples[j][i] = truth
        }
        flip = flip/2
    }
    return samples
}

func (sample Sample) Probability(p float64) float64 {
    sampleP := 1.0
    successes := 0
    for i := 0; i < len(sample); i++ {
        if sample[i] == true {
            successes += 1
        }
    }
    for i := 0; i < successes; i++ {
        sampleP *= p
    }
    for i := successes; i < len(sample); i++ {
        sampleP *= (1-p)
    }
    return sampleP
}


func display(imgData [][]bool) {
    scale := 7
    width := len(imgData)
    height := len(imgData[0])
    upLeft := image.Point{0, 0}
    lowRight := image.Point{width*scale, height*scale}
    img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
    for x := 0; x < width; x++ {
        for xoff := 0; xoff < scale; xoff++ {
            for y := 0; y < height; y++ {
                for yoff := 0; yoff < scale; yoff++ {
                    c := color.Black
                    if imgData[x][y] == true {
                        c = color.White
                    }
                    img.Set(x*scale+xoff, y*scale+yoff, c)
                }
            }
        }
    }
    f, _ := os.Create("image.png")
    png.Encode(f, img)
}

func max(a, b int) int {
    if a > b {
        return a
    }
    return b
}

func main() {
    n := 10
    samples := getSamples(n)
    counts := make([]int, n+1)
    for _, sample := range samples {
        nSuccesses := 0
        for _, trial := range sample {
            if trial {
                nSuccesses += 1
            }
        }
        counts[nSuccesses] += 1
    }
    width := len(counts)
    wScale := 1
    for width < 100 {
        width *= 2
        wScale *= 2
    }
    imgData := make([][]bool, width)
    top := PowInt(2, n)
    height := 100
    for i := 0; i < len(imgData); i++ {
        imgData[i] = make([]bool, height)
        if i == 0 || i == len(imgData)-1 {
            for j := 0; j < height; j++ {
                imgData[i][j] = true
            }
        }
        imgData[i][0] = true
        imgData[i][height-1] = true
    }
    for x, count := range counts {
        percent := float64(count)/float64(top)
        upTo := int(math.Floor(percent*float64(height)))
        for y := 0; y < upTo; y++ {
            imgData[x*wScale][height-1-y] = true
        }
    }
    display(imgData)
    fmt.Println("done")
}
