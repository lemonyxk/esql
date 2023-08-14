/**
* @program: esql
*
* @description:
*
* @author: lemo
*
* @create: 2023-08-14 21:37
**/

package main

import (
	"github.com/lemonyxk/esql"
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {

	var sql = "select * from a where (((id = 1 and name = -1) and addr = china) or b in ('a',    'b',1) or a is not null and c between 10 and 900)"

	log.Println(esql.Convert(sql))

}
