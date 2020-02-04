package main

import (
	"fmt"
	"github.com/isbm/clumad"
)

func main() {
	op := clumad.NewSysOp()
	op.GetSaltOps().SetSaltConfigDir("/etc/salt")
	fmt.Println(op.GetSaltOps().GetConfOpion("minion", "master").(string))
	fmt.Println(op.GetSaltOps().GetConfDOption("minion", "master"))
}
