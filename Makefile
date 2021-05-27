
rebuild-contracts:
	rm -rf bindings/ solidity/
	TARGET=build ./scripts/build_contracts.sh
