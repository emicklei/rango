package main
import "fmt"
import . "bitbucket.org/emicklei/musigo"

func nop(v interface{}){}
func main() {
fmt.Print("")
ParseNote("C") // dummy
d := ParseSequence("C D E F E D E F E D E D E G")
nop(d)
e := PitchBy(1).Transform(d)
nop(e)
fmt.Printf("%v",e)

}
