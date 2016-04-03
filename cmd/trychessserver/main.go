package main

//import "github.com/ArchieT/3manchess/server"
import "github.com/ArchieT/3manchess/server/mojsql"
import "github.com/ArchieT/3manchess/multi"
import "os"
import "log"
import "flag"
import "github.com/coreos/pkg/flagutil"

func main() {
	var mmm mojsql.MojSQL
	flags := flag.NewFlagSet("trychessserver", flag.ExitOnError)
	u := flags.String("u", "", "database username")
	p := flags.String("p", "", "database password")
	d := flags.String("d", "", "database name")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TRYCHESS")
	log.Println(u, p, d)
	log.Println(mmm.Initialize(*u, *p, *d))
	sss := mmm.Interface()
	mul := multi.Multi{sss}
	mul.Run()
}
