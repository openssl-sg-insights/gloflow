/*
GloFlow application and media management/publishing platform
Copyright (C) 2021 Ivan Trajkovic

This program is free software; you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation; either version 2 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program; if not, write to the Free Software
Foundation, Inc., 51 Franklin St, Fifth Floor, Boston, MA  02110-1301  USA
*/

package gf_eth_contract

import (
	// "os"
	"fmt"
	"testing"
	"context"
	"github.com/gloflow/gloflow/go/gf_web3/gf_eth_core"
	// eth_common "github.com/ethereum/go-ethereum/common"
	// "github.com/stretchr/testify/assert"
	// "github.com/gloflow/gloflow/go/gf_core"
	// "github.com/davecgh/go-spew/spew"
)

//---------------------------------------------------
func Test__contract(pTest *testing.T) {

	fmt.Println("TEST__CONTRACT ==============================================")
	
	// block_int := 4634748
	// host_str := os.Getenv("GF_TEST_WORKER_INSPECTOR_HOST")
	// worker_inspector__port_int := 9000
	ctx := context.Background()
	runtime, metrics, err := gf_eth_core.TgetRuntime()
	if err != nil {
		pTest.FailNow()
	}


	abis_map     := t__get_abis()
	erc20_gf_abi := abis_map["erc20"]

	gf_err := Eth_contract__enrich(erc20_gf_abi,
		ctx,
		metrics,
		runtime)
	if gf_err != nil {
		pTest.Fail()
	}

	code_hex_str := "608060405234801561001057600080fd5b5060405160208061021783398101604090815290516000818155338152600160205291909120556101d1806100466000396000f3006080604052600436106100565763ffffffff7c010000000000000000000000000000000000000000000000000000000060003504166318160ddd811461005b57806370a0823114610082578063a9059cbb146100b0575b600080fd5b34801561006757600080fd5b506100706100f5565b60408051918252519081900360200190f35b34801561008e57600080fd5b5061007073ffffffffffffffffffffffffffffffffffffffff600435166100fb565b3480156100bc57600080fd5b506100e173ffffffffffffffffffffffffffffffffffffffff60043516602435610123565b604080519115158252519081900360200190f35b60005490565b73ffffffffffffffffffffffffffffffffffffffff1660009081526001602052604090205490565b600073ffffffffffffffffffffffffffffffffffffffff8316151561014757600080fd5b3360009081526001602052604090205482111561016357600080fd5b503360009081526001602081905260408083208054859003905573ffffffffffffffffffffffffffffffffffffffff85168352909120805483019055929150505600a165627a7a72305820a5d999f4459642872a29be93a490575d345e40fc91a7cccb2cf29c88bcdaf3be0029"
	
	opcodes_lst, gf_err := Eth_contract__get_opcodes(code_hex_str, runtime)
	if gf_err != nil {
		pTest.Fatal()
	}

	// spew.Dump(opcodes_lst)

	for _, gf_opcode := range opcodes_lst {
		fmt.Printf("%s - %s\n", gf_opcode.Addr_hex_str, gf_opcode.Op_and_args_str)
	}
}