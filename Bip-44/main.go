package main

import "fmt"

func main() {
	mnemonic1 := []string{"Arsham", "fox", "Mr.Bean", "number#1", "TooMuch", "Random", "All The same", "Everything", "No Matter What", "nothing", "SoSo", "mm"}

	//creating our first bitcoin account:
	GenerateAddress(mnemonic1, 0, 0, 0, 0)
	fmt.Println("----------------------------------")
	//maybe create an etherium key on another account:
	GenerateAddress(mnemonic1, 60, 0, 1, 0)
	fmt.Println("----------------------------------\n Checking if deterministic: \n----------------------------------")
	GenerateAddress(mnemonic1, 0, 0, 0, 0)
}

func GenerateAddress(mnemonic []string, coinIndex int, changeIndex int, accountIndex int, addressIndex int) string {

	seed := DeriveSeedFromMnemonic(mnemonic, "")
	masterKey := CreateMasterPrivate(seed)

	purpose := DeriveHardened(masterKey, 44)
	coin := DeriveHardened(purpose, uint32(coinIndex))
	account := DeriveHardened(coin, uint32(accountIndex))
	change := DeriveNormalPrivate(account, uint32(changeIndex))
	addressKey := DeriveNormalPrivate(change, uint32(addressIndex))

	pub := PrivateToPublic(addressKey)

	fmt.Printf("Private: %x\n", addressKey.Key.Serialize())
	fmt.Printf("Public: %x\n", pub.Key.SerializeCompressed())

	return string(addressKey.Key.Serialize())

}
