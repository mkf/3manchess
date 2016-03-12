package main

//import "github.com/ArchieT/3manchess/server"
import "github.com/ArchieT/3manchess/server/mojsql"
import "github.com/ArchieT/3manchess/multi"
import "fmt"
import "log"

func main() {
	var mmm mojsql.MojSQL
	var a, b, c string
	fmt.Scanf("%s %s %s", &a, &b, &c)
	log.Println(a, b, c)
	log.Println(mmm.Initialize(a, b, c))
	mul := multi.Multi{&mmm}
	mul.Run()
}
