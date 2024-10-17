package main

import (
	"github.com/amrshaban2005/go-blockchain/utils"
)

func main() {
	myAddress := utils.GetHost()
	utils.FindNeighbors(myAddress, 3333, 0, 3, 3333, 3336)
	//fmt.Println(utils.FindNeighbors("127.0.0.1", 3333, 0, 3, 3333, 3336))
}
